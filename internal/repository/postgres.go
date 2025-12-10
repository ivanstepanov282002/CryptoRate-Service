package repository

import (
	"CryptoRate-Service/internal/models"
	"database/sql"
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