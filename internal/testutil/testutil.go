package testutil

import (
    "cryptorate-service/internal/models"
    "testing"
    "time"
)

// TestTime возвращает фиксированное время для тестов
func TestTime() time.Time {
    return time.Date(2024, 1, 15, 14, 30, 0, 0, time.UTC)
}

// TestCurrency создает тестовую валюту
func TestCurrency() models.Currency {
    return models.Currency{
        ID:           1,
        NameCurrency: "bitcoin",
        DisplayName:  "Bitcoin",
        Symbol:       "BTC",
    }
}

// TestExchangeRate создает тестовый курс
func TestExchangeRate() models.ExchangeRate {
    return models.ExchangeRate{
        ID:         1,
        CurrencyID: 1,
        Price:      45000.50,
        RecordedAt: TestTime(),
    }
}

// AssertNoError проверяет что ошибки нет
func AssertNoError(t *testing.T, err error) {
    t.Helper()
    if err != nil {
        t.Fatalf("Unexpected error: %v", err)
    }
}

// AssertError проверяет что ошибка есть
func AssertError(t *testing.T, err error) {
    t.Helper()
    if err == nil {
        t.Fatal("Expected error, got nil")
    }
}

// AssertEqual проверяет равенство
func AssertEqual[T comparable](t *testing.T, got, want T) {
    t.Helper()
    if got != want {
        t.Errorf("Got %v, want %v", got, want)
    }
}
