package testutil

import (
    "testing"
    "time"
)

func TestTestTime(t *testing.T) {
    result := TestTime()
    expected := time.Date(2024, 1, 15, 14, 30, 0, 0, time.UTC)
    
    if result != expected {
        t.Errorf("TestTime() = %v, want %v", result, expected)
    }
}

func TestTestCurrency(t *testing.T) {
    currency := TestCurrency()
    
    if currency.ID != 1 {
        t.Errorf("Expected ID 1, got %d", currency.ID)
    }
    if currency.NameCurrency != "bitcoin" {
        t.Errorf("Expected NameCurrency 'bitcoin', got '%s'", currency.NameCurrency)
    }
    if currency.DisplayName != "Bitcoin" {
        t.Errorf("Expected DisplayName 'Bitcoin', got '%s'", currency.DisplayName)
    }
    if currency.Symbol != "BTC" {
        t.Errorf("Expected Symbol 'BTC', got '%s'", currency.Symbol)
    }
}

func TestTestExchangeRate(t *testing.T) {
    rate := TestExchangeRate()
    
    if rate.ID != 1 {
        t.Errorf("Expected ID 1, got %d", rate.ID)
    }
    if rate.CurrencyID != 1 {
        t.Errorf("Expected CurrencyID 1, got %d", rate.CurrencyID)
    }
    if rate.Price != 45000.50 {
        t.Errorf("Expected Price 45000.50, got %f", rate.Price)
    }
}

func TestAssertNoError(t *testing.T) {
    // This should not panic
    AssertNoError(t, nil)
}

func TestAssertError(t *testing.T) {
    // This should not panic
    AssertError(t, assertErrorTestError{})
}

func TestAssertEqual(t *testing.T) {
    // This should not panic
    AssertEqual(t, 5, 5)
    
    // Test with strings
    AssertEqual(t, "test", "test")
    
    // Test with floats
    AssertEqual(t, 3.14, 3.14)
}

// Helper type for testing AssertError
type assertErrorTestError struct{}

func (assertErrorTestError) Error() string {
    return "test error"
}