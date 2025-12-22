package api

import (
    "net/http"
    "net/http/httptest"
    "testing"
    "time"
)

func TestCoinGeckoClient_GetPrices(t *testing.T) {
    // Создаем тестовый сервер
    testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Проверяем параметры запроса
        if r.URL.Query().Get("ids") != "bitcoin,ethereum" {
            t.Errorf("Expected ids=bitcoin,ethereum, got %s", r.URL.Query().Get("ids"))
        }
        if r.URL.Query().Get("vs_currencies") != "usd" {
            t.Errorf("Expected vs_currencies=usd, got %s", r.URL.Query().Get("vs_currencies"))
        }

        // Возвращаем тестовый ответ
        w.Header().Set("Content-Type", "application/json")
        w.Write([]byte(`{
            "bitcoin": {"usd": 45000.50},
            "ethereum": {"usd": 2500.75}
        }`))
    }))
    defer testServer.Close()

    // Создаем клиент с тестовым URL
    client := &CoinGeckoClient{
        baseURL: testServer.URL,
        client:  testServer.Client(),
    }

    // Вызываем метод
    prices, err := client.GetPrices([]string{"bitcoin", "ethereum"})
    if err != nil {
        t.Fatalf("GetPrices failed: %v", err)
    }

    // Проверяем результат
    if len(prices) != 2 {
        t.Errorf("Expected 2 prices, got %d", len(prices))
    }

    if btc, ok := prices["bitcoin"]; !ok || btc.USD != 45000.50 {
        t.Errorf("Bitcoin price incorrect. Got %+v", btc)
    }

    if eth, ok := prices["ethereum"]; !ok || eth.USD != 2500.75 {
        t.Errorf("Ethereum price incorrect. Got %+v", eth)
    }
}

func TestCoinGeckoClient_GetPrices_Error(t *testing.T) {
    // Сервер возвращает ошибку
    testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusInternalServerError)
    }))
    defer testServer.Close()

    client := &CoinGeckoClient{
        baseURL: testServer.URL,
        client:  testServer.Client(),
    }

    _, err := client.GetPrices([]string{"bitcoin"})
    if err == nil {
        t.Error("Expected error for failed request")
    }
}

func TestCoinGeckoClient_GetPrices_InvalidJSON(t *testing.T) {
    // Сервер возвращает некорректный JSON
    testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Write([]byte(`invalid json`))
    }))
    defer testServer.Close()

    client := &CoinGeckoClient{
        baseURL: testServer.URL,
        client:  testServer.Client(),
    }

    _, err := client.GetPrices([]string{"bitcoin"})
    if err == nil {
        t.Error("Expected error for invalid JSON")
    }
}

func TestNewCoinGeckoClient(t *testing.T) {
    client := NewCoinGeckoClient()

    if client == nil {
        t.Fatal("Expected non-nil client")
    }

    if client.baseURL != "https://api.coingecko.com/api/v3" {
        t.Errorf("Expected base URL 'https://api.coingecko.com/api/v3', got '%s'", client.baseURL)
    }

    if client.client == nil {
        t.Error("Expected non-nil HTTP client")
    }
}

func TestCoinGeckoClient_GetPrices_Timeout(t *testing.T) {
    // Создаем сервер, который не отвечает (имитация таймаута)
    testServer := httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Никогда не отправляем ответ, чтобы вызвать таймаут
        time.Sleep(100 * time.Millisecond) // Небольшая задержка
    }))
    testServer.Start()
    defer testServer.Close()

    // Клиент с очень маленьким таймаутом
    client := &CoinGeckoClient{
        baseURL: testServer.URL,
        client: &http.Client{
            Timeout: 1 * time.Millisecond, // Очень маленький таймаут
        },
    }

    _, err := client.GetPrices([]string{"bitcoin"})
    if err == nil {
        t.Error("Expected timeout error")
    }
}

// Табличные тесты для разных валют
func TestCoinGeckoClient_GetPrices_Table(t *testing.T) {
    testCases := []struct {
        name     string
        coinIDs  []string
        response string
        wantErr  bool
        desc     string
    }{
        {
            name:    "Single currency",
            coinIDs: []string{"bitcoin"},
            response: `{"bitcoin": {"usd": 45000.50}}`,
            wantErr: false,
            desc:    "Should parse single currency",
        },
        {
            name:    "Multiple currencies",
            coinIDs: []string{"bitcoin", "ethereum", "solana"},
            response: `{"bitcoin": {"usd": 45000.50}, "ethereum": {"usd": 2500.75}, "solana": {"usd": 100.25}}`,
            wantErr: false,
            desc:    "Should parse multiple currencies",
        },
        {
            name:    "Empty response",
            coinIDs: []string{"bitcoin"},
            response: `{}`,
            wantErr: false,
            desc:    "Should handle empty response",
        },
        {
            name:    "Invalid coin ID",
            coinIDs: []string{"invalidcoin"},
            response: `{"invalidcoin": {"usd": 0}}`,
            wantErr: false,
            desc:    "Should handle invalid coin ID",
        },
    }

    for _, tc := range testCases {
        t.Run(tc.name, func(t *testing.T) {
            testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
                w.Header().Set("Content-Type", "application/json")
                w.Write([]byte(tc.response))
            }))
            defer testServer.Close()

            client := &CoinGeckoClient{
                baseURL: testServer.URL,
                client:  testServer.Client(),
            }

            prices, err := client.GetPrices(tc.coinIDs)

            if tc.wantErr && err == nil {
                t.Errorf("%s: expected error, got nil", tc.desc)
            }
            if !tc.wantErr && err != nil {
                t.Errorf("%s: unexpected error: %v", tc.desc, err)
            }

            // Проверяем количество возвращенных курсов
            if err == nil {
                expectedCount := 0
                if tc.response != "{}" {
                    // Грубая проверка - считаем фигурные скобки
                    expectedCount = len(tc.coinIDs)
                }

                if len(prices) != expectedCount && tc.response != "{}" {
                    t.Errorf("%s: expected %d prices, got %d", tc.desc, expectedCount, len(prices))
                }
            }
        })
    }
}

// Бенчмарк тест
func BenchmarkCoinGeckoClient_GetPrices(b *testing.B) {
    testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "application/json")
        w.Write([]byte(`{
            "bitcoin": {"usd": 45000.50},
            "ethereum": {"usd": 2500.75},
            "solana": {"usd": 100.25}
        }`))
    }))
    defer testServer.Close()

    client := &CoinGeckoClient{
        baseURL: testServer.URL,
        client:  testServer.Client(),
    }

    coinIDs := []string{"bitcoin", "ethereum", "solana"}

    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        _, err := client.GetPrices(coinIDs)
        if err != nil {
            b.Fatalf("GetPrices failed: %v", err)
        }
    }
}
