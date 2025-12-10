package models

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