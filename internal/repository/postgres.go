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
