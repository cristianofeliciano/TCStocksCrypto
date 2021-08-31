package crypto

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/luannevesbtc/TCStocksCrypto/model"
	"github.com/luannevesbtc/TCStocksCrypto/store"
	"github.com/nats-io/nats.go"
	"github.com/tradersclub/TCUtils/cache"
	"github.com/tradersclub/TCUtils/tcerr"
)

const (
	URL_GET_MARKETS = "https://api.coingecko.com/api/v3/coins/markets?vs_currency=%s&order=market_cap_desc&per_page=%d&page=%d&sparkline=false&price_change_percentage='1h24h7d'"
	MAX_PER_PAGE    = 250
	BRL             = "BRL"
	USD             = "USD"
	BTC             = "BTC"
)

type getCryptoMarkets struct {
	Err  error
	Id   string
	Data []model.Market
}

// App interface de item para implementação
type App interface {
	GetCryptoMarkets(ctx context.Context) ([]model.Market, error)
}

// NewApp cria uma nova instancia do serviço de exemplo item
func NewApp(stores *store.Container, nc *nats.Conn, cache cache.Cache) App {
	return &appImpl{
		stores: stores,
		nc:     nc,
		cache:  cache,
	}
}

type appImpl struct {
	stores    *store.Container
	startedAt time.Time
	version   string
	nc        *nats.Conn
	cache     cache.Cache
}

func (s *appImpl) GetCryptoMarkets(ctx context.Context) ([]model.Market, error) {
	response := make([]model.Market, 0)
	i := 1
	for {
		brlMarkets, err := getMarketRequest(BRL, i)
		if err != nil {
			return nil, err
		}
		btcMarkets, err := getMarketRequest(BTC, i)
		if err != nil {
			return nil, err
		}
		usdMarkets, err := getMarketRequest(USD, i)
		if err != nil {
			return nil, err
		}

		if len(usdMarkets) == 0 {
			break
		}

		response = append(response, btcMarkets...)
		response = append(response, brlMarkets...)
		response = append(response, usdMarkets...)
		i++
	}
	
	return response, nil
}

func getMarketRequest(currency string, page int) ([]model.Market, error) {
	jsonData := make([]model.Market, 0)

	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf(URL_GET_MARKETS, currency, MAX_PER_PAGE, page), nil)
	if err != nil {
		return nil, tcerr.NewError(http.StatusInternalServerError, "erro ao realizar o get dos markets", currency)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, tcerr.NewError(http.StatusInternalServerError, "erro ao realizar o get dos markets", currency)
	}

	decoder := json.NewDecoder(resp.Body)
	errResponse := decoder.Decode(&jsonData)
	if errResponse != nil {
		return nil, tcerr.NewError(http.StatusInternalServerError, "erro ao realizar o parser dos markets", currency)
	}

	return jsonData, nil
}
