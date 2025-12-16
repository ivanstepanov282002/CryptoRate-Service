package models

import "time"

type Currency struct {
	ID           int    `json:"id"`
	NameCurrency string `json:"name_currency"`
}

type ExchangeRate struct {
	ID         int     `json:"id"`
	CurrencyID int     `json:"currency_id"`
	Price      float64 `json:"price"`
}

type CoinGeckoResponse map[string]struct {
	USD float64 `json:"usd"`
}

type CurrencyRateView struct {
	NameCurrency string    `json:"name_currency"`
	Price        float64   `json:"price"`
	RecordedAt   time.Time `json:"recorded_at"`
}
