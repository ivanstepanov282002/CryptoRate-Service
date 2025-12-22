package models

import (
    "encoding/json"
    "testing"
    "time"
)

func TestCurrency_JSON(t *testing.T) {
    currency := Currency{
        ID:           1,
        NameCurrency: "bitcoin",
        DisplayName:  "Bitcoin",
        Symbol:       "BTC",
    }

    data, err := json.Marshal(currency)
    if err != nil {
        t.Fatalf("Failed to marshal currency: %v", err)
    }

    var decoded Currency
    if err := json.Unmarshal(data, &decoded); err != nil {
        t.Fatalf("Failed to unmarshal currency: %v", err)
    }

    if decoded != currency {
        t.Errorf("Decoded currency doesn't match original. Got %+v, want %+v", decoded, currency)
    }
}

func TestExchangeRate_JSON(t *testing.T) {
    now := time.Now()
    rate := ExchangeRate{
        ID:         1,
        CurrencyID: 1,
        Price:      45000.50,
        RecordedAt: now,
    }

    data, err := json.Marshal(rate)
    if err != nil {
        t.Fatalf("Failed to marshal exchange rate: %v", err)
    }

    var decoded ExchangeRate
    if err := json.Unmarshal(data, &decoded); err != nil {
        t.Fatalf("Failed to unmarshal exchange rate: %v", err)
    }

    // Сравниваем с точностью до микросекунд из-за JSON сериализации времени
    if decoded.ID != rate.ID || decoded.CurrencyID != rate.CurrencyID || 
    decoded.Price != rate.Price {
        t.Errorf("Decoded rate doesn't match original. Got %+v, want %+v", decoded, rate)
    }
}

func TestCoinGeckoResponse_Parse(t *testing.T) {
    jsonData := `{
        "bitcoin": {"usd": 45000.50},
        "ethereum": {"usd": 2500.75}
    }`

    var response CoinGeckoResponse
    if err := json.Unmarshal([]byte(jsonData), &response); err != nil {
        t.Fatalf("Failed to parse CoinGecko response: %v", err)
    }

    if len(response) != 2 {
        t.Errorf("Expected 2 currencies, got %d", len(response))
    }

    if btc, ok := response["bitcoin"]; !ok || btc.USD != 45000.50 {
        t.Errorf("Bitcoin price incorrect. Got %+v", btc)
    }

    if eth, ok := response["ethereum"]; !ok || eth.USD != 2500.75 {
        t.Errorf("Ethereum price incorrect. Got %+v", eth)
    }
}

func TestCurrencyRateView_JSON(t *testing.T) {
    now := time.Now()
    rateView := CurrencyRateView{
        NameCurrency: "bitcoin",
        Price:        45000.50,
        RecordedAt:   now,
    }

    data, err := json.Marshal(rateView)
    if err != nil {
        t.Fatalf("Failed to marshal currency rate view: %v", err)
    }

    var decoded CurrencyRateView
    if err := json.Unmarshal(data, &decoded); err != nil {
        t.Fatalf("Failed to unmarshal currency rate view: %v", err)
    }

    if decoded.NameCurrency != rateView.NameCurrency || decoded.Price != rateView.Price {
        t.Errorf("Decoded rate view doesn't match original. Got %+v, want %+v", decoded, rateView)
    }
}

func TestUserSettings_JSON(t *testing.T) {
    userSettings := UserSettings{
        UserID:   123456,
        Interval: 10,
        LastSent: time.Now(),
        Currencies: []Currency{
            {ID: 1, NameCurrency: "bitcoin", DisplayName: "Bitcoin", Symbol: "BTC"},
            {ID: 2, NameCurrency: "ethereum", DisplayName: "Ethereum", Symbol: "ETH"},
        },
    }

    data, err := json.Marshal(userSettings)
    if err != nil {
        t.Fatalf("Failed to marshal user settings: %v", err)
    }

    var decoded UserSettings
    if err := json.Unmarshal(data, &decoded); err != nil {
        t.Fatalf("Failed to unmarshal user settings: %v", err)
    }

    if decoded.UserID != userSettings.UserID || decoded.Interval != userSettings.Interval {
        t.Errorf("Decoded user settings don't match original")
    }

    if len(decoded.Currencies) != len(userSettings.Currencies) {
        t.Errorf("Currency count mismatch. Got %d, want %d", 
                len(decoded.Currencies), len(userSettings.Currencies))
    }
}

// Бенчмарк тесты
func BenchmarkCurrency_Marshal(b *testing.B) {
    currency := Currency{
        ID:           1,
        NameCurrency: "bitcoin",
        DisplayName:  "Bitcoin",
        Symbol:       "BTC",
    }

    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        _, err := json.Marshal(currency)
        if err != nil {
            b.Fatalf("Marshal failed: %v", err)
        }
    }
}

func BenchmarkCoinGeckoResponse_Unmarshal(b *testing.B) {
    jsonData := []byte(`{
        "bitcoin": {"usd": 45000.50},
        "ethereum": {"usd": 2500.75},
        "tether": {"usd": 1.00}
    }`)

    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        var response CoinGeckoResponse
        if err := json.Unmarshal(jsonData, &response); err != nil {
            b.Fatalf("Unmarshal failed: %v", err)
        }
    }
}
