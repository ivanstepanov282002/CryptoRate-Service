package repository

import (
    "cryptorate-service/internal/models"
    "database/sql"
    "testing"
    "time"

    _ "github.com/lib/pq"
    "github.com/DATA-DOG/go-sqlmock"
)

func TestNewRepository(t *testing.T) {
    db, _, err := sqlmock.New()
    if err != nil {
        t.Fatalf("Failed to create sqlmock: %v", err)
    }
    defer db.Close()

    repo := NewRepository(db)
    if repo.db != db {
        t.Error("Repository should use the provided database connection")
    }
}

func TestRepository_SaveRate(t *testing.T) {
    db, mock, err := sqlmock.New()
    if err != nil {
        t.Fatalf("Failed to create sqlmock: %v", err)
    }
    defer db.Close()

    repo := NewRepository(db)

    // Тест успешного сохранения
    rate := models.ExchangeRate{
        CurrencyID: 1,
        Price:      100.50,
    }

    mock.ExpectExec(`INSERT INTO exchange_rate \(currency_id, price\)`).
        WithArgs(rate.CurrencyID, rate.Price).
        WillReturnResult(sqlmock.NewResult(1, 1))

    err = repo.SaveRate(rate)
    if err != nil {
        t.Errorf("SaveRate failed: %v", err)
    }

    if err := mock.ExpectationsWereMet(); err != nil {
        t.Errorf("Unfulfilled expectations: %v", err)
    }
}

func TestRepository_GetCurrencyID(t *testing.T) {
    db, mock, err := sqlmock.New()
    if err != nil {
        t.Fatalf("Failed to create sqlmock: %v", err)
    }
    defer db.Close()

    repo := NewRepository(db)

    // Тест успешного получения ID
    expectedID := 1
    currencyName := "bitcoin"

    rows := sqlmock.NewRows([]string{"id"}).
        AddRow(expectedID)

    mock.ExpectQuery("SELECT id FROM currency WHERE LOWER\\(name_currency\\) = LOWER\\(\\$1\\)").
        WithArgs(currencyName).
        WillReturnRows(rows)

    id, err := repo.GetCurrencyID(currencyName)
    if err != nil {
        t.Errorf("GetCurrencyID failed: %v", err)
    }

    if id != expectedID {
        t.Errorf("Expected ID %d, got %d", expectedID, id)
    }

    if err := mock.ExpectationsWereMet(); err != nil {
        t.Errorf("Unfulfilled expectations: %v", err)
    }

    // Тест ошибки при отсутствии валюты
    db2, mock2, err := sqlmock.New()
    if err != nil {
        t.Fatalf("Failed to create sqlmock: %v", err)
    }
    defer db2.Close()

    repo2 := NewRepository(db2)

    mock2.ExpectQuery("SELECT id FROM currency WHERE LOWER\\(name_currency\\) = LOWER\\(\\$1\\)").
        WithArgs("nonexistent").
        WillReturnError(sql.ErrNoRows)

    _, err = repo2.GetCurrencyID("nonexistent")
    if err == nil {
        t.Error("Expected error for non-existent currency")
    }

    if err := mock2.ExpectationsWereMet(); err != nil {
        t.Errorf("Unfulfilled expectations: %v", err)
    }
}

func TestRepository_GetCurrencyIDBySymbol(t *testing.T) {
    db, mock, err := sqlmock.New()
    if err != nil {
        t.Fatalf("Failed to create sqlmock: %v", err)
    }
    defer db.Close()

    repo := NewRepository(db)

    // Тест успешного получения ID по символу
    expectedID := 1
    symbol := "BTC"

    rows := sqlmock.NewRows([]string{"id"}).
        AddRow(expectedID)

    mock.ExpectQuery("SELECT id FROM currency WHERE LOWER\\(symbol\\) = LOWER\\(\\$1\\)").
        WithArgs(symbol).
        WillReturnRows(rows)

    id, err := repo.GetCurrencyIDBySymbol(symbol)
    if err != nil {
        t.Errorf("GetCurrencyIDBySymbol failed: %v", err)
    }

    if id != expectedID {
        t.Errorf("Expected ID %d, got %d", expectedID, id)
    }

    if err := mock.ExpectationsWereMet(); err != nil {
        t.Errorf("Unfulfilled expectations: %v", err)
    }
}

func TestRepository_GetCurrencySymbol(t *testing.T) {
    db, mock, err := sqlmock.New()
    if err != nil {
        t.Fatalf("Failed to create sqlmock: %v", err)
    }
    defer db.Close()

    repo := NewRepository(db)

    // Тест успешного получения символа
    expectedSymbol := "BTC"
    currencyID := 1

    rows := sqlmock.NewRows([]string{"symbol"}).
        AddRow(expectedSymbol)

    mock.ExpectQuery("SELECT symbol FROM currency WHERE id = \\$1").
        WithArgs(currencyID).
        WillReturnRows(rows)

    symbol, err := repo.GetCurrencySymbol(currencyID)
    if err != nil {
        t.Errorf("GetCurrencySymbol failed: %v", err)
    }

    if symbol != expectedSymbol {
        t.Errorf("Expected symbol %s, got %s", expectedSymbol, symbol)
    }

    if err := mock.ExpectationsWereMet(); err != nil {
        t.Errorf("Unfulfilled expectations: %v", err)
    }
}

func TestRepository_GetCurrencySymbolByID(t *testing.T) {
    db, mock, err := sqlmock.New()
    if err != nil {
        t.Fatalf("Failed to create sqlmock: %v", err)
    }
    defer db.Close()

    repo := NewRepository(db)

    // Тест успешного получения символа по ID
    expectedSymbol := "BTC"
    currencyID := 1

    rows := sqlmock.NewRows([]string{"symbol"}).
        AddRow(expectedSymbol)

    mock.ExpectQuery("SELECT symbol FROM currency WHERE id = \\$1").
        WithArgs(currencyID).
        WillReturnRows(rows)

    symbol, err := repo.GetCurrencySymbolByID(currencyID)
    if err != nil {
        t.Errorf("GetCurrencySymbolByID failed: %v", err)
    }

    if symbol != expectedSymbol {
        t.Errorf("Expected symbol %s, got %s", expectedSymbol, symbol)
    }

    if err := mock.ExpectationsWereMet(); err != nil {
        t.Errorf("Unfulfilled expectations: %v", err)
    }
}

func TestRepository_GetCurrencyDisplayName(t *testing.T) {
    db, mock, err := sqlmock.New()
    if err != nil {
        t.Fatalf("Failed to create sqlmock: %v", err)
    }
    defer db.Close()

    repo := NewRepository(db)

    // Тест успешного получения отображаемого имени
    expectedName := "Bitcoin"
    currencyID := 1

    rows := sqlmock.NewRows([]string{"display_name"}).
        AddRow(expectedName)

    mock.ExpectQuery("SELECT display_name FROM currency WHERE id = \\$1").
        WithArgs(currencyID).
        WillReturnRows(rows)

    name, err := repo.GetCurrencyDisplayName(currencyID)
    if err != nil {
        t.Errorf("GetCurrencyDisplayName failed: %v", err)
    }

    if name != expectedName {
        t.Errorf("Expected name %s, got %s", expectedName, name)
    }

    if err := mock.ExpectationsWereMet(); err != nil {
        t.Errorf("Unfulfilled expectations: %v", err)
    }
}

func TestRepository_GetCurrencyRate(t *testing.T) {
    db, mock, err := sqlmock.New()
    if err != nil {
        t.Fatalf("Failed to create sqlmock: %v", err)
    }
    defer db.Close()

    repo := NewRepository(db)

    // Тест успешного получения курса
    currencyID := 1
    expectedRate := models.ExchangeRate{
        ID:         1,
        CurrencyID: currencyID,
        Price:      45000.50,
        RecordedAt: time.Now(),
    }

    rows := sqlmock.NewRows([]string{"id", "currency_id", "price", "recorded_at"}).
        AddRow(expectedRate.ID, expectedRate.CurrencyID, expectedRate.Price, expectedRate.RecordedAt)

    mock.ExpectQuery(`SELECT id, currency_id, price, recorded_at FROM Exchange_rate WHERE currency_id = \$1 ORDER BY recorded_at DESC LIMIT 1`).
        WithArgs(currencyID).
        WillReturnRows(rows)

    rate, err := repo.GetCurrencyRate(currencyID)
    if err != nil {
        t.Errorf("GetCurrencyRate failed: %v", err)
    }

    if rate.CurrencyID != expectedRate.CurrencyID {
        t.Errorf("Expected currency ID %d, got %d", expectedRate.CurrencyID, rate.CurrencyID)
    }

    if rate.Price != expectedRate.Price {
        t.Errorf("Expected price %.2f, got %.2f", expectedRate.Price, rate.Price)
    }

    if err := mock.ExpectationsWereMet(); err != nil {
        t.Errorf("Unfulfilled expectations: %v", err)
    }
}

func TestRepository_GetDailyMinMax(t *testing.T) {
    db, mock, err := sqlmock.New()
    if err != nil {
        t.Fatalf("Failed to create sqlmock: %v", err)
    }
    defer db.Close()

    repo := NewRepository(db)

    // Тест успешного получения дневного мин/макс
    currencyID := 1
    expectedMin := 44000.00
    expectedMax := 46000.00

    rows := sqlmock.NewRows([]string{"min", "max"}).
        AddRow(expectedMin, expectedMax)

    mock.ExpectQuery(`SELECT MIN\(price\), MAX\(price\) FROM Exchange_rate WHERE currency_id = \$1`).
        WithArgs(currencyID).
        WillReturnRows(rows)

    min, max, err := repo.GetDailyMinMax(currencyID)
    if err != nil {
        t.Errorf("GetDailyMinMax failed: %v", err)
    }

    if min != expectedMin {
        t.Errorf("Expected min %.2f, got %.2f", expectedMin, min)
    }

    if max != expectedMax {
        t.Errorf("Expected max %.2f, got %.2f", expectedMax, max)
    }

    if err := mock.ExpectationsWereMet(); err != nil {
        t.Errorf("Unfulfilled expectations: %v", err)
    }
}

func TestRepository_GetHourlyChange(t *testing.T) {
    db, mock, err := sqlmock.New()
    if err != nil {
        t.Fatalf("Failed to create sqlmock: %v", err)
    }
    defer db.Close()

    repo := NewRepository(db)

    // Тест успешного получения часового изменения
    currencyID := 1
    currentPrice := 45000.50
    priceHourAgo := 44000.00
    expectedChange := (currentPrice - priceHourAgo) / priceHourAgo * 100

    // Сначала текущая цена
    rows1 := sqlmock.NewRows([]string{"price"}).
        AddRow(currentPrice)
    mock.ExpectQuery(`SELECT price FROM Exchange_rate WHERE currency_id = \$1 ORDER BY recorded_at DESC LIMIT 1`).
        WithArgs(currencyID).
        WillReturnRows(rows1)

    // Затем цена час назад
    rows2 := sqlmock.NewRows([]string{"price"}).
        AddRow(priceHourAgo)
    mock.ExpectQuery(`SELECT price FROM Exchange_rate WHERE currency_id = \$1 AND recorded_at <= NOW\(\) - INTERVAL '1 hour' ORDER BY recorded_at DESC LIMIT 1`).
        WithArgs(currencyID).
        WillReturnRows(rows2)

    change, err := repo.GetHourlyChange(currencyID)
    if err != nil {
        t.Errorf("GetHourlyChange failed: %v", err)
    }

    // Допускаем небольшое отклонение из-за округления
    if change != expectedChange {
        t.Errorf("Expected change %.2f, got %.2f", expectedChange, change)
    }

    if err := mock.ExpectationsWereMet(); err != nil {
        t.Errorf("Unfulfilled expectations: %v", err)
    }
}

func TestRepository_GetLatestRates(t *testing.T) {
    db, mock, err := sqlmock.New()
    if err != nil {
        t.Fatalf("Failed to create sqlmock: %v", err)
    }
    defer db.Close()

    repo := NewRepository(db)

    // Тест успешного получения последних курсов
    rows := sqlmock.NewRows([]string{"name_currency", "price", "recorded_at", "currency_id"}).
        AddRow("bitcoin", 45000.50, time.Now(), 1).
        AddRow("ethereum", 3000.25, time.Now(), 2)

    mock.ExpectQuery(`SELECT DISTINCT ON \(c\.name_currency\) c\.name_currency, e\.price, e\.recorded_at, c\.id as currency_id FROM currency c JOIN exchange_rate e ON c\.id = e\.currency_id ORDER BY c\.name_currency, e\.recorded_at DESC`).
        WillReturnRows(rows)

    rates, err := repo.GetLatestRates()
    if err != nil {
        t.Errorf("GetLatestRates failed: %v", err)
    }

    if len(rates) != 2 {
        t.Errorf("Expected 2 rates, got %d", len(rates))
    }

    if err := mock.ExpectationsWereMet(); err != nil {
        t.Errorf("Unfulfilled expectations: %v", err)
    }
}

func TestRepository_GetAllCurrencies(t *testing.T) {
    db, mock, err := sqlmock.New()
    if err != nil {
        t.Fatalf("Failed to create sqlmock: %v", err)
    }
    defer db.Close()

    repo := NewRepository(db)

    // Тест успешного получения всех валют
    rows := sqlmock.NewRows([]string{"id", "name_currency", "display_name", "symbol"}).
        AddRow(1, "bitcoin", "Bitcoin", "BTC").
        AddRow(2, "ethereum", "Ethereum", "ETH")

    mock.ExpectQuery(`SELECT id, name_currency, display_name, symbol FROM currency ORDER BY id`).
        WillReturnRows(rows)

    currencies, err := repo.GetAllCurrencies()
    if err != nil {
        t.Errorf("GetAllCurrencies failed: %v", err)
    }

    if len(currencies) != 2 {
        t.Errorf("Expected 2 currencies, got %d", len(currencies))
    }

    if currencies[0].NameCurrency != "bitcoin" {
        t.Errorf("Expected first currency to be 'bitcoin', got '%s'", currencies[0].NameCurrency)
    }

    if err := mock.ExpectationsWereMet(); err != nil {
        t.Errorf("Unfulfilled expectations: %v", err)
    }
}

func TestRepository_Ping(t *testing.T) {
    db, mock, err := sqlmock.New()
    if err != nil {
        t.Fatalf("Failed to create sqlmock: %v", err)
    }
    defer db.Close()

    repo := NewRepository(db)

    mock.ExpectPing()

    err = repo.Ping()
    if err != nil {
        t.Errorf("Ping failed: %v", err)
    }

    if err := mock.ExpectationsWereMet(); err != nil {
        t.Errorf("Unfulfilled expectations: %v", err)
    }
}

// Табличные тесты для разных сценариев GetCurrencyID
func TestRepository_GetCurrencyID_Table(t *testing.T) {
    testCases := []struct {
        name     string
        setupMock func(sqlmock.Sqlmock, string)
        input    string
        wantErr  bool
        desc     string
    }{
        {
            "Bitcoin lowercase",
            func(mock sqlmock.Sqlmock, currency string) {
                rows := sqlmock.NewRows([]string{"id"}).AddRow(1)
                mock.ExpectQuery("SELECT id FROM currency WHERE LOWER\\(name_currency\\) = LOWER\\(\\$1\\)").
                    WithArgs(currency).
                    WillReturnRows(rows)
            },
            "bitcoin",
            false,
            "Should find bitcoin",
        },
        {
            "Bitcoin uppercase",
            func(mock sqlmock.Sqlmock, currency string) {
                rows := sqlmock.NewRows([]string{"id"}).AddRow(1)
                mock.ExpectQuery("SELECT id FROM currency WHERE LOWER\\(name_currency\\) = LOWER\\(\\$1\\)").
                    WithArgs(currency).
                    WillReturnRows(rows)
            },
            "BITCOIN",
            false,
            "Should be case insensitive",
        },
        {
            "Nonexistent currency",
            func(mock sqlmock.Sqlmock, currency string) {
                mock.ExpectQuery("SELECT id FROM currency WHERE LOWER\\(name_currency\\) = LOWER\\(\\$1\\)").
                    WithArgs(currency).
                    WillReturnError(sql.ErrNoRows)
            },
            "nonexistentcoin",
            true,
            "Should error on nonexistent currency",
        },
        {
            "Empty string",
            func(mock sqlmock.Sqlmock, currency string) {
                mock.ExpectQuery("SELECT id FROM currency WHERE LOWER\\(name_currency\\) = LOWER\\(\\$1\\)").
                    WithArgs(currency).
                    WillReturnError(sql.ErrNoRows)
            },
            "",
            true,
            "Should error on empty string",
        },
    }

    for _, tc := range testCases {
        t.Run(tc.name, func(t *testing.T) {
            db, mock, err := sqlmock.New()
            if err != nil {
                t.Fatalf("Failed to create sqlmock: %v", err)
            }
            defer db.Close()

            repo := NewRepository(db)
            tc.setupMock(mock, tc.input)

            _, err = repo.GetCurrencyID(tc.input)

            if tc.wantErr && err == nil {
                t.Errorf("%s: expected error, got nil", tc.desc)
            }
            if !tc.wantErr && err != nil {
                t.Errorf("%s: unexpected error: %v", tc.desc, err)
            }

            if err := mock.ExpectationsWereMet(); err != nil {
                t.Errorf("Unfulfilled expectations: %v", err)
            }
        })
    }
}
