package repository

import (
	"cryptorate-service/internal/models"
	"database/sql"
	"fmt"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

// SaveRate saves the currency exchange rate in the database
func (r *Repository) SaveRate(rate models.ExchangeRate) error {
	_, err := r.db.Exec(`
        INSERT INTO exchange_rate (currency_id, price) 
        VALUES ($1, $2)`,
		rate.CurrencyID, rate.Price)
	return err
}

// GetCurrencyID возвращает ID валюты по её имени
func (r *Repository) GetCurrencyID(name string) (int, error) {
	var id int
	err := r.db.QueryRow("SELECT id FROM currency WHERE LOWER(name_currency) = LOWER($1)", name).Scan(&id)
	return id, err
}

// GetLatestRates вовращает послдний курс каждой вылюты из БД
func (r *Repository) GetLatestRates() ([]models.CurrencyRateView, error) {
	// SQL запрос: для каждой валюты берём самую свежую запись
	query := `
        SELECT DISTINCT ON (c.name_currency)
            c.name_currency,
            e.price,
            e.recorded_at
        FROM currency c
        JOIN exchange_rate e ON c.id = e.currency_id
        ORDER BY c.name_currency, e.recorded_at DESC`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("query failed: %w", err)
	}
	defer rows.Close()

	//Создаем слайс, который хранит в себе строки базы данных (название валюты, цена, время получения курса)
	var rates []models.CurrencyRateView
	for rows.Next() {
		//1 строка из базы данных
		var rate models.CurrencyRateView
		err := rows.Scan(&rate.NameCurrency, &rate.Price, &rate.RecordedAt)
		if err != nil {
			return nil, fmt.Errorf("scan failed: %w", err)
		}
		//Добавляем в слайс 1 строку с валютой
		rates = append(rates, rate)
	}

	return rates, nil
}

// GetCurrencySymbol возвращает символ валюты по её ID
func (r *Repository) GetCurrencySymbol(currencyID int) (string, error) {
    var symbol string
    err := r.db.QueryRow("SELECT symbol FROM currency WHERE id = $1", currencyID).Scan(&symbol)
    return symbol, err
}

// GetCurrencyDisplayName возвращает отображаемое имя валюты
func (r *Repository) GetCurrencyDisplayName(currencyID int) (string, error) {
    var displayName string
    err := r.db.QueryRow("SELECT display_name FROM currency WHERE id = $1", currencyID).Scan(&displayName)
    return displayName, err
}

// GetAllCurrencies возвращает все доступные валюты
func (r *Repository) GetAllCurrencies() ([]models.Currency, error) {
    query := "SELECT id, name_currency, display_name, symbol FROM currency ORDER BY id"
    
    rows, err := r.db.Query(query)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var currencies []models.Currency
    for rows.Next() {
        var currency models.Currency
        err := rows.Scan(&currency.ID, &currency.NameCurrency, 
                        &currency.DisplayName, &currency.Symbol)
        if err != nil {
            return nil, err
        }
        currencies = append(currencies, currency)
    }

    return currencies, nil
}

// GetCurrencyIDBySymbol возвращает ID валюты по символу (BTC, ETH)
func (r *Repository) GetCurrencyIDBySymbol(symbol string) (int, error) {
    var id int
    err := r.db.QueryRow(
        "SELECT id FROM currency WHERE LOWER(symbol) = LOWER($1)", 
        symbol,
    ).Scan(&id)
    return id, err
}

// GetCurrencySymbolByID возвращает символ валюты по ID
func (r *Repository) GetCurrencySymbolByID(currencyID int) (string, error) {
    var symbol string
    err := r.db.QueryRow(
        "SELECT symbol FROM currency WHERE id = $1", 
        currencyID,
    ).Scan(&symbol)
    return symbol, err
}

// EnsureUser создает пользователя, если его нет
func (r *Repository) EnsureUser(userID int64, userName string) error {
    _, err := r.db.Exec(`
        INSERT INTO Users (user_id, user_name) 
        VALUES ($1, $2)
        ON CONFLICT (user_id) DO UPDATE SET user_name = $2
    `, userID, userName)
    return err
}

// SetUserInterval устанавливает интервал автоотправки
func (r *Repository) SetUserInterval(userID int64, interval int) error {
    // Убедимся что пользователь есть
    r.EnsureUser(userID, "")
    
    tx, err := r.db.Begin()
    if err != nil {
        return err
    }

    // Обновляем или добавляем настройки
    _, err = tx.Exec(`
        INSERT INTO Settings (user_id, time_interval, last_sent) 
        VALUES ($1, $2, NOW())
        ON CONFLICT (user_id) 
        DO UPDATE SET time_interval = $2, last_sent = NOW()
    `, userID, interval)
    if err != nil {
        tx.Rollback()
        return err
    }

    // Активируем все валюты для пользователя
    _, err = tx.Exec(`
        INSERT INTO Currency_settings (user_id, currency_id, is_active)
        SELECT $1, id, true
        FROM Currency
        ON CONFLICT (user_id, currency_id) 
        DO UPDATE SET is_active = true
    `, userID)
    if err != nil {
        tx.Rollback()
        return err
    }

    return tx.Commit()
}

// StopAuto отключает автоотправку
func (r *Repository) StopAuto(userID int64) error {
    _, err := r.db.Exec(`
        UPDATE Settings 
        SET time_interval = NULL
        WHERE user_id = $1
    `, userID)
    return err
}

// GetSubscribedUsers возвращает пользователей с активными подписками
func (r *Repository) GetSubscribedUsers() ([]models.UserSettings, error) {
    query := `
        SELECT s.user_id, s.time_interval, s.last_sent, 
               c.id, c.name_currency, c.display_name, c.symbol
        FROM Settings s
        JOIN Currency_settings cs ON s.user_id = cs.user_id AND cs.is_active = true
        JOIN Currency c ON cs.currency_id = c.id
        WHERE s.time_interval > 0
        ORDER BY s.user_id`

    rows, err := r.db.Query(query)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var users []models.UserSettings
    var currentUser *models.UserSettings
    var lastUserID int64

    for rows.Next() {
        var userID int64
        var interval int
        var lastSent time.Time
        var currencyID int
        var nameCurrency, displayName, symbol string

        err := rows.Scan(&userID, &interval, &lastSent, 
                        &currencyID, &nameCurrency, &displayName, &symbol)
        if err != nil {
            return nil, err
        }

        if lastUserID != userID {
            users = append(users, models.UserSettings{
                UserID:     userID,
                Interval:   interval,
                LastSent:   lastSent,
                Currencies: []models.Currency{{ID: currencyID, NameCurrency: nameCurrency, 
                                             DisplayName: displayName, Symbol: symbol}},
            })
            currentUser = &users[len(users)-1]
            lastUserID = userID
        } else {
            currentUser.Currencies = append(currentUser.Currencies, 
                models.Currency{ID: currencyID, NameCurrency: nameCurrency, 
                               DisplayName: displayName, Symbol: symbol})
        }
    }

    return users, nil
}

// UpdateLastSent обновляет время последней отправки
func (r *Repository) UpdateLastSent(userID int64) error {
    _, err := r.db.Exec(`
        UPDATE Settings 
        SET last_sent = NOW()
        WHERE user_id = $1
    `, userID)
    return err
}

// GetCurrencyRate возвращает последний курс для валюты
func (r *Repository) GetCurrencyRate(currencyID int) (models.ExchangeRate, error) {
    var rate models.ExchangeRate
    err := r.db.QueryRow(`
        SELECT id, currency_id, price, recorded_at
        FROM Exchange_rate
        WHERE currency_id = $1
        ORDER BY recorded_at DESC
        LIMIT 1`, currencyID).Scan(&rate.ID, &rate.CurrencyID, &rate.Price, &rate.RecordedAt)
    return rate, err
}

// GetDailyMinMax возвращает минимальную и максимальную цену за сегодня
func (r *Repository) GetDailyMinMax(currencyID int) (min, max float64, err error) {
    query := `
        SELECT MIN(price), MAX(price)
        FROM Exchange_rate
        WHERE currency_id = $1 
          AND recorded_at >= CURRENT_DATE
          AND recorded_at < CURRENT_DATE + INTERVAL '1 day'`
    
    err = r.db.QueryRow(query, currencyID).Scan(&min, &max)
    return
}

// GetHourlyChange возвращает изменение цены за последний час в процентах
func (r *Repository) GetHourlyChange(currencyID int) (change float64, err error) {
    // Текущая цена
    var currentPrice float64
    err = r.db.QueryRow(`
        SELECT price 
        FROM Exchange_rate 
        WHERE currency_id = $1 
        ORDER BY recorded_at DESC 
        LIMIT 1`, currencyID).Scan(&currentPrice)
    if err != nil {
        return 0, err
    }

    // Цена час назад
    var priceHourAgo float64
    err = r.db.QueryRow(`
        SELECT price 
        FROM Exchange_rate 
        WHERE currency_id = $1 
          AND recorded_at <= NOW() - INTERVAL '1 hour'
        ORDER BY recorded_at DESC 
        LIMIT 1`, currencyID).Scan(&priceHourAgo)
    if err != nil {
        // Если нет записи час назад, возвращаем 0
        return 0, nil
    }

    if priceHourAgo == 0 {
        return 0, nil
    }

    change = (currentPrice - priceHourAgo) / priceHourAgo * 100
    return change, nil
}

// DB возвращает соединение с БД (для health check)
func (r *Repository) DB() *sql.DB {
    return r.db
}

// Ping проверяет соединение с БД
func (r *Repository) Ping() error {
    return r.db.Ping()
}