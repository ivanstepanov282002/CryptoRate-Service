package rest

import (
	"cryptorate-service/internal/repository"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/mux"
)

type Response struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
	Meta    *Meta       `json:"meta,omitempty"`
}

type Meta struct {
	Timestamp string `json:"timestamp"`
	Version   string `json:"version"`
}

type RateResponse struct {
	Currency     string    `json:"currency"`
	Symbol       string    `json:"symbol"`
	DisplayName  string    `json:"display_name"`
	Price        float64   `json:"price"`
	UpdatedAt    time.Time `json:"updated_at"`
	DailyMin     float64   `json:"daily_min,omitempty"`
	DailyMax     float64   `json:"daily_max,omitempty"`
	HourlyChange float64   `json:"hourly_change,omitempty"`
}

type StatsResponse struct {
	Currency     string    `json:"currency"`
	Symbol       string    `json:"symbol"`
	DisplayName  string    `json:"display_name"` 
	Current      float64   `json:"current"`
	DailyMin     float64   `json:"daily_min"`
	DailyMax     float64   `json:"daily_max"`
	HourlyChange float64   `json:"hourly_change"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type CurrencyResponse struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	DisplayName string `json:"display_name"`
	Symbol      string `json:"symbol"`
}

type Handler struct {
	repo *repository.Repository
}

func NewHandler(repo *repository.Repository) *Handler {
	return &Handler{repo: repo}
}

// GetRates возвращает все курсы
func (h *Handler) GetRates(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	rates, err := h.repo.GetLatestRates()
	if err != nil {
		sendError(w, "Failed to get rates", http.StatusInternalServerError)
		return
	}

	if len(rates) == 0 {
		sendError(w, "No rates found", http.StatusNotFound)
		return
	}

	response := make([]RateResponse, len(rates))
	for i, rate := range rates {
		currencyID, err := h.repo.GetCurrencyID(rate.NameCurrency)
		if err != nil {
			continue
		}

		symbol, _ := h.repo.GetCurrencySymbolByID(currencyID)
		displayName, _ := h.repo.GetCurrencyDisplayName(currencyID)
		min, max, _ := h.repo.GetDailyMinMax(currencyID)
		change, _ := h.repo.GetHourlyChange(currencyID)

		response[i] = RateResponse{
			Currency:     rate.NameCurrency,
			Symbol:       symbol,
			DisplayName:  displayName,
			Price:        rate.Price,
			UpdatedAt:    rate.RecordedAt,
			DailyMin:     min,
			DailyMax:     max,
			HourlyChange: change,
		}
	}

	sendJSON(w, Response{
		Success: true,
		Data:    response,
		Meta:    &Meta{Timestamp: time.Now().Format(time.RFC3339), Version: "1.0"},
	})
}

// GetRate возвращает курс конкретной валюты
func (h *Handler) GetRate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	currencyName := strings.ToLower(vars["currency"])

	// Пробуем найти по символу или имени
	currencyID, err := h.repo.GetCurrencyIDBySymbol(currencyName)
	if err != nil {
		currencyID, err = h.repo.GetCurrencyID(currencyName)
	}

	if err != nil {
		sendError(w, "Currency not found", http.StatusNotFound)
		return
	}

	rate, err := h.repo.GetCurrencyRate(currencyID)
	if err != nil {
		sendError(w, "Rate not found", http.StatusNotFound)
		return
	}

	symbol, _ := h.repo.GetCurrencySymbolByID(currencyID)
	displayName, _ := h.repo.GetCurrencyDisplayName(currencyID)
	min, max, _ := h.repo.GetDailyMinMax(currencyID)
	change, _ := h.repo.GetHourlyChange(currencyID)

	response := RateResponse{
		Currency:     currencyName,
		Symbol:       symbol,
		DisplayName:  displayName,
		Price:        rate.Price,
		UpdatedAt:    rate.RecordedAt,
		DailyMin:     min,
		DailyMax:     max,
		HourlyChange: change,
	}

	sendJSON(w, Response{
		Success: true,
		Data:    response,
		Meta:    &Meta{Timestamp: time.Now().Format(time.RFC3339), Version: "1.0"},
	})
}

// GetStats возвращает расширенную статистику
func (h *Handler) GetStats(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	currencyName := strings.ToLower(vars["currency"])

	currencyID, err := h.repo.GetCurrencyIDBySymbol(currencyName)
	if err != nil {
		currencyID, err = h.repo.GetCurrencyID(currencyName)
	}

	if err != nil {
		sendError(w, "Currency not found", http.StatusNotFound)
		return
	}

	rate, err := h.repo.GetCurrencyRate(currencyID)
	if err != nil {
		sendError(w, "Rate not found", http.StatusNotFound)
		return
	}

	symbol, _ := h.repo.GetCurrencySymbolByID(currencyID)
	displayName, _ := h.repo.GetCurrencyDisplayName(currencyID)
	min, max, _ := h.repo.GetDailyMinMax(currencyID)
	change, _ := h.repo.GetHourlyChange(currencyID)

	response := StatsResponse{
		Currency:     currencyName,
		Symbol:       symbol,
		DisplayName:  displayName,
		Current:      rate.Price,
		DailyMin:     min,
		DailyMax:     max,
		HourlyChange: change,
		UpdatedAt:    rate.RecordedAt,
	}

	sendJSON(w, Response{
		Success: true,
		Data:    response,
		Meta:    &Meta{Timestamp: time.Now().Format(time.RFC3339), Version: "1.0"},
	})
}

// GetCurrencies возвращает список всех валют
func (h *Handler) GetCurrencies(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	currencies, err := h.repo.GetAllCurrencies()
	if err != nil {
		sendError(w, "Failed to get currencies", http.StatusInternalServerError)
		return
	}

	response := make([]CurrencyResponse, len(currencies))
	for i, currency := range currencies {
		response[i] = CurrencyResponse{
			ID:          currency.ID,
			Name:        currency.NameCurrency,
			DisplayName: currency.DisplayName,
			Symbol:      currency.Symbol,
		}
	}

	sendJSON(w, Response{
		Success: true,
		Data:    response,
		Meta:    &Meta{Timestamp: time.Now().Format(time.RFC3339), Version: "1.0"},
	})
}

// HealthCheck проверяет состояние сервиса
func (h *Handler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Проверяем соединение с БД
	dbErr := h.repo.Ping()

	health := map[string]interface{}{
		"status":    "healthy",
		"timestamp": time.Now().Unix(),
		"service":   "crypto-rates-api",
		"version":   "1.0.0",
		"database":  "connected",
		"uptime":    time.Since(startTime).String(),
	}

	if dbErr != nil {
		health["status"] = "unhealthy"
		health["database"] = "disconnected"
		health["error"] = dbErr.Error()
	}

	sendJSON(w, Response{
		Success: true,
		Data:    health,
		Meta:    &Meta{Timestamp: time.Now().Format(time.RFC3339), Version: "1.0"},
	})
}

// Вспомогательные методы
func sendJSON(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

func sendError(w http.ResponseWriter, message string, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(Response{
		Success: false,
		Error:   message,
		Meta:    &Meta{Timestamp: time.Now().Format(time.RFC3339), Version: "1.0"},
	})
}

var startTime = time.Now()
