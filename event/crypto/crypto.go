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

const TCGET_CRYPTO_MARKETS = "tcget_crypto_markets"

type getMarkets struct {
	Err  error
	Id   string
	Data []model.Market
}

// Register group health check
func Register(apps *app.Container, conn *nats.Conn) {
	e := &event{
		apps: apps,
		nc:   conn,
	}

	e.nc.Subscribe(TCGET_CRYPTO_MARKETS, e.getMarkets)
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

	//b, _ := json.Marshal(getMarketsResponse)
}
