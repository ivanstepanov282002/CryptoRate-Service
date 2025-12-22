package rest

import (
    "cryptorate-service/internal/models"
    "encoding/json"
    "fmt"
    "net/http"
    "net/http/httptest"
    "testing"
    "time"

    "github.com/gorilla/mux"
)

// MockRepository для тестирования
type MockRepository struct {
    rates      []models.CurrencyRateView
    currencies []models.Currency
    err        error
}

func (m *MockRepository) GetLatestRates() ([]models.CurrencyRateView, error) {
    return m.rates, m.err
}

func (m *MockRepository) GetAllCurrencies() ([]models.Currency, error) {
    return m.currencies, m.err
}

func (m *MockRepository) GetCurrencyID(name string) (int, error) {
    if name == "bitcoin" || name == "btc" {
        return 1, m.err
    } else if name == "ethereum" || name == "eth" {
        return 2, m.err
    }
    return 0, fmt.Errorf("currency not found: %s", name)
}

func (m *MockRepository) GetCurrencyIDBySymbol(symbol string) (int, error) {
    if symbol == "BTC" || symbol == "btc" {
        return 1, m.err
    } else if symbol == "ETH" || symbol == "eth" {
        return 2, m.err
    }
    return 0, fmt.Errorf("symbol not found: %s", symbol)
}

func (m *MockRepository) GetCurrencyRate(currencyID int) (models.ExchangeRate, error) {
    return models.ExchangeRate{
        ID:         1,
        CurrencyID: currencyID,
        Price:      45000.50,
        RecordedAt: time.Now(),
    }, m.err
}

func (m *MockRepository) GetDailyMinMax(currencyID int) (min, max float64, err error) {
    return 44500.00, 45500.75, m.err
}

func (m *MockRepository) GetHourlyChange(currencyID int) (change float64, err error) {
    return 1.25, m.err
}

func (m *MockRepository) GetCurrencySymbolByID(currencyID int) (string, error) {
    if currencyID == 1 {
        return "BTC", m.err
    } else if currencyID == 2 {
        return "ETH", m.err
    }
    return "", fmt.Errorf("currency ID not found: %d", currencyID)
}

func (m *MockRepository) GetCurrencyDisplayName(currencyID int) (string, error) {
    if currencyID == 1 {
        return "Bitcoin", m.err
    } else if currencyID == 2 {
        return "Ethereum", m.err
    }
    return "", fmt.Errorf("currency ID not found: %d", currencyID)
}

func (m *MockRepository) Ping() error {
    return m.err
}

func (m *MockRepository) DB() interface{} {
    return nil
}

func (m *MockRepository) GetCurrencySymbol(currencyID int) (string, error) {
    if currencyID == 1 {
        return "BTC", m.err
    } else if currencyID == 2 {
        return "ETH", m.err
    }
    return "", fmt.Errorf("currency ID not found: %d", currencyID)
}

func TestHandler_GetRates(t *testing.T) {
    // Подготовка мок данных
    mockRates := []models.CurrencyRateView{
        {
            NameCurrency: "bitcoin",
            Price:        45000.50,
            RecordedAt:   time.Now(),
        },
    }

    repo := &MockRepository{rates: mockRates}
    handler := NewHandler(repo)

    // Создание запроса
    req := httptest.NewRequest("GET", "/api/v1/rates", nil)
    w := httptest.NewRecorder()

    // Выполнение запроса
    handler.GetRates(w, req)

    // Проверка статуса
    if w.Code != http.StatusOK {
        t.Errorf("Expected status 200, got %d", w.Code)
    }

    // Проверка JSON ответа
    var response Response
    if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
        t.Fatalf("Failed to parse response: %v", err)
    }

    if !response.Success {
        t.Error("Expected success to be true")
    }

    // Проверка данных
    ratesData, ok := response.Data.([]interface{})
    if !ok || len(ratesData) == 0 {
        t.Error("Expected rates data in response")
    }
}

func TestHandler_GetCurrencies(t *testing.T) {
    mockCurrencies := []models.Currency{
        {ID: 1, NameCurrency: "bitcoin", DisplayName: "Bitcoin", Symbol: "BTC"},
        {ID: 2, NameCurrency: "ethereum", DisplayName: "Ethereum", Symbol: "ETH"},
    }

    repo := &MockRepository{currencies: mockCurrencies}
    handler := NewHandler(repo)

    req := httptest.NewRequest("GET", "/api/v1/currencies", nil)
    w := httptest.NewRecorder()

    handler.GetCurrencies(w, req)

    if w.Code != http.StatusOK {
        t.Errorf("Expected status 200, got %d", w.Code)
    }

    var response Response
    if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
        t.Fatalf("Failed to parse response: %v", err)
    }

    if !response.Success {
        t.Error("Expected success to be true")
    }
}

func TestHandler_HealthCheck(t *testing.T) {
    repo := &MockRepository{}
    handler := NewHandler(repo)

    req := httptest.NewRequest("GET", "/api/v1/health", nil)
    w := httptest.NewRecorder()

    handler.HealthCheck(w, req)

    if w.Code != http.StatusOK {
        t.Errorf("Expected status 200, got %d", w.Code)
    }

    var response Response
    if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
        t.Fatalf("Failed to parse response: %v", err)
    }

    if !response.Success {
        t.Error("Expected success to be true")
    }

    // Проверка health данных
    healthData, ok := response.Data.(map[string]interface{})
    if !ok {
        t.Error("Expected health data in response")
    }

    if status, ok := healthData["status"].(string); !ok || status == "" {
        t.Error("Expected status in health response")
    }
}

func TestHandler_GetRate(t *testing.T) {
    // Подготовка мок данных
    repo := &MockRepository{}
    handler := NewHandler(repo)

    // Создание запроса с параметром
    req := httptest.NewRequest("GET", "/api/v1/rate/bitcoin", nil)
    req = mux.SetURLVars(req, map[string]string{"currency": "bitcoin"})
    w := httptest.NewRecorder()

    // Выполнение запроса
    handler.GetRate(w, req)

    // Проверка статуса
    if w.Code != http.StatusOK {
        t.Errorf("Expected status 200, got %d", w.Code)
    }

    // Проверка JSON ответа
    var response Response
    if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
        t.Fatalf("Failed to parse response: %v", err)
    }

    if !response.Success {
        t.Error("Expected success to be true")
    }

    // Проверка данных
    _, ok := response.Data.(map[string]interface{})
    if !ok {
        t.Error("Expected rate data in response")
    }
}

func TestHandler_GetRate_Error(t *testing.T) {
    // Тестирование ошибки при поиске валюты
    repo := &MockRepository{err: fmt.Errorf("currency not found")}
    handler := NewHandler(repo)

    // Создание запроса с параметром
    req := httptest.NewRequest("GET", "/api/v1/rate/nonexistent", nil)
    req = mux.SetURLVars(req, map[string]string{"currency": "nonexistent"})
    w := httptest.NewRecorder()

    // Выполнение запроса
    handler.GetRate(w, req)

    // Проверка статуса - должен быть 404
    if w.Code != http.StatusNotFound {
        t.Errorf("Expected status 404, got %d", w.Code)
    }

    // Проверка JSON ответа
    var response Response
    if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
        t.Fatalf("Failed to parse response: %v", err)
    }

    if response.Success {
        t.Error("Expected success to be false")
    }

    if response.Error == "" {
        t.Error("Expected error message in response")
    }
}

func TestHandler_GetRate_NotFound(t *testing.T) {
    // Тестирование случая, когда валюта не найдена
    repo := &MockRepository{}
    handler := NewHandler(repo)

    // Создание запроса с несуществующей валютой
    req := httptest.NewRequest("GET", "/api/v1/rate/nonexistent", nil)
    req = mux.SetURLVars(req, map[string]string{"currency": "nonexistent"})
    w := httptest.NewRecorder()

    // Выполнение запроса
    handler.GetRate(w, req)

    // Проверка статуса - должен быть 404
    if w.Code != http.StatusNotFound {
        t.Errorf("Expected status 404, got %d", w.Code)
    }

    // Проверка JSON ответа
    var response Response
    if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
        t.Fatalf("Failed to parse response: %v", err)
    }

    if response.Success {
        t.Error("Expected success to be false")
    }

    if response.Error == "" {
        t.Error("Expected error message in response")
    }
}

func TestHandler_GetStats(t *testing.T) {
    // Подготовка мок данных
    repo := &MockRepository{}
    handler := NewHandler(repo)

    // Создание запроса с параметром
    req := httptest.NewRequest("GET", "/api/v1/stats/bitcoin", nil)
    req = mux.SetURLVars(req, map[string]string{"currency": "bitcoin"})
    w := httptest.NewRecorder()

    // Выполнение запроса
    handler.GetStats(w, req)

    // Проверка статуса
    if w.Code != http.StatusOK {
        t.Errorf("Expected status 200, got %d", w.Code)
    }

    // Проверка JSON ответа
    var response Response
    if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
        t.Fatalf("Failed to parse response: %v", err)
    }

    if !response.Success {
        t.Error("Expected success to be true")
    }

    // Проверка данных
    _, ok := response.Data.(map[string]interface{})
    if !ok {
        t.Error("Expected stats data in response")
    }
}

func TestHandler_GetStats_Error(t *testing.T) {
    // Тестирование ошибки при поиске статистики
    repo := &MockRepository{err: fmt.Errorf("currency not found")}
    handler := NewHandler(repo)

    // Создание запроса с параметром
    req := httptest.NewRequest("GET", "/api/v1/stats/nonexistent", nil)
    req = mux.SetURLVars(req, map[string]string{"currency": "nonexistent"})
    w := httptest.NewRecorder()

    // Выполнение запроса
    handler.GetStats(w, req)

    // Проверка статуса - должен быть 404
    if w.Code != http.StatusNotFound {
        t.Errorf("Expected status 404, got %d", w.Code)
    }

    // Проверка JSON ответа
    var response Response
    if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
        t.Fatalf("Failed to parse response: %v", err)
    }

    if response.Success {
        t.Error("Expected success to be false")
    }

    if response.Error == "" {
        t.Error("Expected error message in response")
    }
}

func TestHandler_GetStats_NotFound(t *testing.T) {
    // Тестирование случая, когда статистика не найдена
    repo := &MockRepository{}
    handler := NewHandler(repo)

    // Создание запроса с несуществующей валютой
    req := httptest.NewRequest("GET", "/api/v1/stats/nonexistent", nil)
    req = mux.SetURLVars(req, map[string]string{"currency": "nonexistent"})
    w := httptest.NewRecorder()

    // Выполнение запроса
    handler.GetStats(w, req)

    // Проверка статуса - должен быть 404
    if w.Code != http.StatusNotFound {
        t.Errorf("Expected status 404, got %d", w.Code)
    }

    // Проверка JSON ответа
    var response Response
    if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
        t.Fatalf("Failed to parse response: %v", err)
    }

    if response.Success {
        t.Error("Expected success to be false")
    }

    if response.Error == "" {
        t.Error("Expected error message in response")
    }
}

func TestHandler_GetRates_Empty(t *testing.T) {
    // Тестирование случая, когда нет курсов валют
    repo := &MockRepository{rates: []models.CurrencyRateView{}}
    handler := NewHandler(repo)

    // Создание запроса
    req := httptest.NewRequest("GET", "/api/v1/rates", nil)
    w := httptest.NewRecorder()

    // Выполнение запроса
    handler.GetRates(w, req)

    // Проверка статуса - должен быть 404
    if w.Code != http.StatusNotFound {
        t.Errorf("Expected status 404, got %d", w.Code)
    }

    // Проверка JSON ответа
    var response Response
    if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
        t.Fatalf("Failed to parse response: %v", err)
    }

    if response.Success {
        t.Error("Expected success to be false")
    }

    if response.Error == "" {
        t.Error("Expected error message in response")
    }
}

func TestHandler_GetRates_Error(t *testing.T) {
    // Тестирование ошибки при получении курсов
    repo := &MockRepository{err: fmt.Errorf("database error")}
    handler := NewHandler(repo)

    // Создание запроса
    req := httptest.NewRequest("GET", "/api/v1/rates", nil)
    w := httptest.NewRecorder()

    // Выполнение запроса
    handler.GetRates(w, req)

    // Проверка статуса - должен быть 500
    if w.Code != http.StatusInternalServerError {
        t.Errorf("Expected status 500, got %d", w.Code)
    }

    // Проверка JSON ответа
    var response Response
    if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
        t.Fatalf("Failed to parse response: %v", err)
    }

    if response.Success {
        t.Error("Expected success to be false")
    }

    if response.Error == "" {
        t.Error("Expected error message in response")
    }
}

func TestHandler_GetCurrencies_Error(t *testing.T) {
    // Тестирование ошибки при получении списка валют
    repo := &MockRepository{err: fmt.Errorf("database error")}
    handler := NewHandler(repo)

    // Создание запроса
    req := httptest.NewRequest("GET", "/api/v1/currencies", nil)
    w := httptest.NewRecorder()

    // Выполнение запроса
    handler.GetCurrencies(w, req)

    // Проверка статуса - должен быть 500
    if w.Code != http.StatusInternalServerError {
        t.Errorf("Expected status 500, got %d", w.Code)
    }

    // Проверка JSON ответа
    var response Response
    if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
        t.Fatalf("Failed to parse response: %v", err)
    }

    if response.Success {
        t.Error("Expected success to be false")
    }

    if response.Error == "" {
        t.Error("Expected error message in response")
    }
}

func TestHandler_HealthCheck_Unhealthy(t *testing.T) {
    // Тестирование случая, когда сервис нездоров
    repo := &MockRepository{err: fmt.Errorf("database connection failed")}
    handler := NewHandler(repo)

    req := httptest.NewRequest("GET", "/api/v1/health", nil)
    w := httptest.NewRecorder()

    handler.HealthCheck(w, req)

    if w.Code != http.StatusOK {
        t.Errorf("Expected status 200, got %d", w.Code)
    }

    var response Response
    if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
        t.Fatalf("Failed to parse response: %v", err)
    }

    if !response.Success {
        t.Error("Expected success to be true")
    }

    // Проверка health данных
    healthData, ok := response.Data.(map[string]interface{})
    if !ok {
        t.Error("Expected health data in response")
    }

    if status, ok := healthData["status"].(string); !ok || status != "unhealthy" {
        t.Error("Expected status to be 'unhealthy' when database is down")
    }
}

// Тестирование вспомогательных функций
func TestSendJSON(t *testing.T) {
    w := httptest.NewRecorder()
    testData := map[string]string{"test": "data"}

    sendJSON(w, testData)

    if w.Code != http.StatusOK {
        t.Errorf("Expected status 200, got %d", w.Code)
    }

    if w.Header().Get("Content-Type") != "application/json" {
        t.Error("Expected Content-Type: application/json")
    }
}

func TestSendError(t *testing.T) {
    w := httptest.NewRecorder()
    message := "Test error message"
    code := http.StatusBadRequest

    sendError(w, message, code)

    if w.Code != code {
        t.Errorf("Expected status %d, got %d", code, w.Code)
    }

    var response Response
    if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
        t.Fatalf("Failed to parse error response: %v", err)
    }

    if response.Success {
        t.Error("Expected success to be false for error response")
    }

    if response.Error != message {
        t.Errorf("Expected error message '%s', got '%s'", message, response.Error)
    }
}

// Бенчмарк тест
func BenchmarkHealthCheck(b *testing.B) {
    repo := &MockRepository{}
    handler := NewHandler(repo)

    req := httptest.NewRequest("GET", "/api/v1/health", nil)

    for i := 0; i < b.N; i++ {
        w := httptest.NewRecorder()
        handler.HealthCheck(w, req)
    }
}