package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

type TonPrice struct {
	PriceUsd              float64 `json:"usd"`
	PriceChangePercent24H float64 `json:"usd_24h_change"`
	MarketCapUsd          float64 `json:"usd_market_cap"`
	Volume24HUsd          float64 `json:"usd_24h_vol"`
}

type TonResponse struct {
	ID                    string  `json:"id"`
	Name                  string  `json:"name"`
	Symbol                string  `json:"symbol"`
	PriceUsd              float64 `json:"price_usd"`
	PriceChangePercent24H float64 `json:"price_change_percent_24h"`
	MarketCapUsd          float64 `json:"market_cap_usd"`
	Volume24HUsd          float64 `json:"volume_24h_usd"`
	LastUpdated           string  `json:"last_updated"`
}

func getTonPriceHandler(c echo.Context) error {
	price, err := getTonPriceFromCoinGecko()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to fetch TON price: "+err.Error())
	}

	response := TonResponse{
		ID:                    "the-open-network",
		Name:                  "The Open Network",
		Symbol:                "TON",
		PriceUsd:              price.PriceUsd,
		PriceChangePercent24H: price.PriceChangePercent24H,
		MarketCapUsd:          price.MarketCapUsd,
		Volume24HUsd:          price.Volume24HUsd,
		LastUpdated:           time.Now().Format(time.RFC3339),
	}

	return c.JSON(http.StatusOK, response)
}

func getTonPriceFromCoinGecko() (*TonPrice, error) {
	url := "https://api.coingecko.com/api/v3/simple/price?ids=the-open-network&vs_currencies=usd&include_24hr_change=true&include_market_cap=true&include_24hr_vol=true"

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("HTTP request failed: %w", err)
	}
	defer resp.Body.Close()

	// Проверяем статус ответа
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status: %d", resp.StatusCode)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("JSON decode failed: %w", err)
	}

	if tonData, exists := result["the-open-network"].(map[string]interface{}); exists {
		price := &TonPrice{}

		// Безопасное извлечение данных с проверкой типов
		if usd, ok := tonData["usd"].(float64); ok {
			price.PriceUsd = usd
		}

		if change, ok := tonData["usd_24h_change"].(float64); ok {
			price.PriceChangePercent24H = change
		}

		if marketCap, ok := tonData["usd_market_cap"].(float64); ok {
			price.MarketCapUsd = marketCap
		}

		if volume, ok := tonData["usd_24h_vol"].(float64); ok {
			price.Volume24HUsd = volume
		}

		return price, nil
	}

	return nil, fmt.Errorf("TON data not found in response")
}

func main() {
	e := echo.New()

	e.GET("/", getTonPriceHandler)

	fmt.Println("🚀 Server started on http://localhost:8000")
	fmt.Println("📊 TON Price API: http://localhost:8000/api/ton/price")
	e.Logger.Fatal(e.Start(":8000"))
}
