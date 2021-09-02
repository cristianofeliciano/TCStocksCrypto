package crypto

import (
	"context"
	"net/http"

	"github.com/luannevesbtc/TCStocksCrypto/app"
	"github.com/luannevesbtc/TCStocksCrypto/model"
	"github.com/nats-io/nats.go"

	"github.com/tradersclub/TCUtils/logger"
	"github.com/tradersclub/TCUtils/tcerr"
)

const (
	TCGET_CRYPTO_MARKETS    = "tcget_crypto_markets"
	TCGET_GLOBAL_INFOS      = "tcget_global_infos"
	TCGET_CRYPTO_CATEGORIES = "tcget_crypto_categories"
	TCGET_CRYPTO_TICKERS    = "tcget_crypto_tickers"
)

type getMarkets struct {
	Err  error
	Id   string
	Data []model.Market
}

type getGlobalInfos struct {
	Err  error
	Id   string
	Data *model.GlobalInfos
}

type getCryptoCategories struct {
	Err  error
	Id   string
	Data []model.CryptoCategories
}

type getCryptoTickers struct {
	Err  error
	Id   string
	Data []model.CryptoTycker
}

// Register group health check
func Register(apps *app.Container, conn *nats.Conn) {
	e := &event{
		apps: apps,
		nc:   conn,
	}

	e.nc.Subscribe(TCGET_CRYPTO_MARKETS, e.getMarkets)
	e.nc.Subscribe(TCGET_GLOBAL_INFOS, e.getGlobalInfos)
	e.nc.Subscribe(TCGET_CRYPTO_CATEGORIES, e.getCryptoCategories)
	e.nc.Subscribe(TCGET_CRYPTO_TICKERS, e.getCryptoTickers)
}

type event struct {
	apps *app.Container
	nc   *nats.Conn
}

func (e *event) getMarkets(msg *nats.Msg) {
	ctx := context.Background()

	var getMarketsResponse getMarkets

	markets, err := e.apps.Crypto.GetCryptoMarkets(ctx)
	if err != nil {
		logger.ErrorContext(ctx, "event.session.get_crypto_markets", err.Error())
		getMarketsResponse.Err = tcerr.NewError(http.StatusInternalServerError, "event.session.get_crypto_markets", err.Error())
	} else {
		getMarketsResponse.Data = markets
	}
}

func (e *event) getGlobalInfos(msg *nats.Msg) {
	ctx := context.Background()

	var getGlobalInfosResponse getGlobalInfos

	infos, err := e.apps.Crypto.GetGlobalInfos(ctx)
	if err != nil {
		logger.ErrorContext(ctx, "event.session.get_global_infos", err.Error())
		getGlobalInfosResponse.Err = tcerr.NewError(http.StatusInternalServerError, "event.session.get_global_infos", err.Error())
	} else {
		getGlobalInfosResponse.Data = infos
	}
}

func (e *event) getCryptoCategories(msg *nats.Msg) {
	ctx := context.Background()

	var getCryptoCategoriesResponse getCryptoCategories

	categories, err := e.apps.Crypto.GetCryptoCategories(ctx)
	if err != nil {
		logger.ErrorContext(ctx, "event.session.get_crypto_categories", err.Error())
		getCryptoCategoriesResponse.Err = tcerr.NewError(http.StatusInternalServerError, "event.session.get_crypto_categories", err.Error())
	} else {
		getCryptoCategoriesResponse.Data = categories
	}
}

func (e *event) getCryptoTickers(msg *nats.Msg) {
	ctx := context.Background()

	var getCryptoTickersResponse getCryptoTickers

	list, err := e.apps.Crypto.GetCryptoList(ctx)
	if err != nil {
		logger.ErrorContext(ctx, "event.session.get_crypto_tickers", err.Error())
		getCryptoTickersResponse.Err = tcerr.NewError(http.StatusInternalServerError, "event.session.get_crypto_tickers", err.Error())
	} else {
		getCryptoTickersResponse.Data = list
	}
}
