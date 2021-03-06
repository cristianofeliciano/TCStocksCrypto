package crypto

import (
	"context"
	"log"
	"net/http"

	"github.com/nats-io/nats.go"
	nmodelCrypto "github.com/tradersclub/TCNatsModel/TCCrypto"
	"github.com/tradersclub/TCStocksCrypto/app"

	nmodel "github.com/tradersclub/TCNatsModel/TCStocksBovespa"
	"github.com/tradersclub/TCUtils/logger"
	"github.com/tradersclub/TCUtils/tcerr"
)

// Register group health check
func Register(apps *app.Container, conn *nats.Conn) {
	e := &event{
		apps: apps,
		nc:   conn,
	}

	e.nc.Subscribe(nmodelCrypto.TCGET_CRYPTO_MARKETS, e.getMarkets)
	e.nc.Subscribe(nmodelCrypto.TCGET_GLOBAL_INFOS, e.getGlobalInfos)
	e.nc.Subscribe(nmodelCrypto.TCGET_CRYPTO_CATEGORIES, e.getCryptoCategories)
	e.nc.Subscribe(nmodelCrypto.TCGET_CRYPTO_TICKERS, e.getCryptoTickers)
}

type event struct {
	apps *app.Container
	nc   *nats.Conn
}

func (e *event) getMarkets(msg *nats.Msg) {
	ctx := context.Background()

	var getMarketsResponse nmodelCrypto.GetMarkets

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

	var getGlobalInfosResponse nmodelCrypto.GetGlobalInfos

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

	var getCryptoCategoriesResponse nmodelCrypto.GetCryptoCategories

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

	var getCryptoTickersResponse nmodelCrypto.GetCryptoTickers

	list, err := e.apps.Crypto.GetCryptoList(ctx)
	if err != nil {
		logger.ErrorContext(ctx, "event.session.get_crypto_tickers", err.Error())
		getCryptoTickersResponse.Err = tcerr.NewError(http.StatusInternalServerError, "event.session.get_crypto_tickers", err.Error())
	} else {
		getCryptoTickersResponse.Data = list
	}

	tickersList := make([]nmodel.Stock, 0)
	for _, ticker := range list {
		tickerStock := nmodel.Stock{
			Nome:               ticker.Name,
			CodigoTicker:       ticker.Symbol,
			Ticker:             ("$" + ticker.Symbol),
			StocksSegmentId:    0,
			Version:            0,
			VersionType:        "M",
			InternalCryptoType: "crypto",
			InternalSymbol:     ticker.Symbol,
			CodigoISINPapel:    ticker.ID,
		}
		tickersList = append(tickersList, tickerStock)
	}

	err = e.insertStock(tickersList)
	if err != nil {
		logger.ErrorContext(ctx, "event.session.insert_stock", err.Error())
		getCryptoTickersResponse.Err = tcerr.NewError(http.StatusInternalServerError, "event.session.insert_stock", err.Error())
	}
}

func (e *event) insertStock(stocks []nmodel.Stock) error {
	if len(stocks) == 0 {
		return nil
	}
	sendStock := []nmodel.Stock{}
	for _, stock := range stocks {
		sendStock = append(sendStock, stock)
		if len(sendStock) == 1000 {
			e.sendNatsInsert(sendStock)
			sendStock = []nmodel.Stock{}
		}
	}
	err := e.sendNatsInsert(sendStock)
	return err
}

func (e *event) sendNatsInsert(stocks []nmodel.Stock) error {
	log.Print("Send insert to NATS: ", len(stocks))
	insertStock := nmodel.InsertStockList{
		Stocks: stocks,
	}
	err := e.nc.Publish(nmodel.NATS_TCSTOCKS_INSERT_STOCKS, insertStock.ToBytes())
	if err != nil {
		return err
	}

	return nil
}
