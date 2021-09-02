package model

import (
	"net/http"
	"time"

	"github.com/tradersclub/TCUtils/tcerr"
)

type Market struct {
	ID                           string    `json:"id"`
	Symbol                       string    `json:"symbol"`
	Name                         string    `json:"name"`
	Image                        string    `json:"image"`
	CurrentPrice                 float64   `json:"current_price"`
	MarketCap                    float64   `json:"market_cap"`
	MarketCapRank                int       `json:"market_cap_rank"`
	FullyDilutedValuation        *float64  `json:"fully_diluted_valuation"`
	TotalVolume                  float64   `json:"total_volume"`
	High24H                      float64   `json:"high_24h"`
	Low24H                       float64   `json:"low_24h"`
	PriceChange24H               float64   `json:"price_change_24h"`
	PriceChangePercentage24H     float64   `json:"price_change_percentage_24h"`
	MarketCapChange24H           float64   `json:"market_cap_change_24h"`
	MarketCapChangePercentage24H float64   `json:"market_cap_change_percentage_24h"`
	CirculatingSupply            float64   `json:"circulating_supply"`
	TotalSupply                  *float64  `json:"total_supply"`
	MaxSupply                    *float64  `json:"max_supply"`
	Ath                          float64   `json:"ath"`
	AthChangePercentage          float64   `json:"ath_change_percentage"`
	AthDate                      time.Time `json:"ath_date"`
	Atl                          float64   `json:"atl"`
	AtlChangePercentage          float64   `json:"atl_change_percentage"`
	AtlDate                      time.Time `json:"atl_date"`
	Roi                          *struct {
		Times      float64 `json:"times"`
		Currency   string  `json:"currency"`
		Percentage float64 `json:"percentage"`
	} `json:"roi"`
	LastUpdated                        time.Time `json:"last_updated"`
	PriceChangePercentage1HInCurrency  float64   `json:"price_change_percentage_1h_in_currency"`
	PriceChangePercentage24HInCurrency float64   `json:"price_change_percentage_24h_in_currency"`
	PriceChangePercentage7DInCurrency  float64   `json:"price_change_percentage_7d_in_currency"`
}

type GlobalInfos struct {
	Data struct {
		ActiveCryptocurrencies          int                `json:"active_cryptocurrencies"`
		UpcomingIcos                    int                `json:"upcoming_icos"`
		OngoingIcos                     int                `json:"ongoing_icos"`
		EndedIcos                       int                `json:"ended_icos"`
		Markets                         int                `json:"markets"`
		TotalMarketCap                  map[string]float64 `json:"total_market_cap"`
		TotalVolume                     map[string]float64 `json:"total_volume"`
		MarketCapPercentage             map[string]float64 `json:"market_cap_percentage"`
		MarketCapChangePercentage24HUsd float64            `json:"market_cap_change_percentage_24h_usd"`
		UpdatedAt                       int                `json:"updated_at"`
	} `json:"data"`
}

type CryptoCategories struct {
	ID                 string    `json:"id"`
	Name               string    `json:"name"`
	MarketCap          float64   `json:"market_cap"`
	MarketCapChange24H float64   `json:"market_cap_change_24h"`
	Volume24H          float64   `json:"volume_24h"`
	UpdatedAt          time.Time `json:"updated_at"`
}

type CryptoTycker struct {
	ID     string `json:"id"`
	Symbol string `json:"symbol"`
	Name   string `json:"name"`
}

// ToMarket converte uma interface{} para *Market
func ToMarket(data interface{}) (*Market, error) {
	value, ok := data.(*Market)
	if !ok {
		return nil, tcerr.New(http.StatusInternalServerError, "não foi possível converter interface{} para *ToMarket", nil)
	}
	return value, nil
}

// ToGlobalInfos converte uma interface{} para *GlobalInfos
func ToGlobalInfos(data interface{}) (*GlobalInfos, error) {
	value, ok := data.(*GlobalInfos)
	if !ok {
		return nil, tcerr.New(http.StatusInternalServerError, "não foi possível converter interface{} para *GlobalInfos", nil)
	}
	return value, nil
}
