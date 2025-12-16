package api

import (
	"cryptorate-service/internal/models"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type CoinGeckoClient struct {
	baseURL string
	client  *http.Client
}

func NewCoinGeckoClient() *CoinGeckoClient {
	return &CoinGeckoClient{
		baseURL: "https://api.coingecko.com/api/v3",
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

//Выполняет запрос курса валют по API, читает ответ, парсит JSON
func (c *CoinGeckoClient) GetPrices(coinIDs []string) (models.CoinGeckoResponse, error) {
	// Формируем URL
	params := url.Values{}
	params.Add("ids", strings.Join(coinIDs, ","))
	params.Add("vs_currencies", "usd")
	url := fmt.Sprintf("%s/simple/price?%s", c.baseURL, params.Encode())

	// Выполняем запрос
	resp, err := c.client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("API request failed: %w", err)
	}
	defer resp.Body.Close()

	// Читаем ответ
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	// Парсим JSON
	var result models.CoinGeckoResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	return result, nil
}
