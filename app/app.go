package app

import (
	"time"

	"github.com/nats-io/nats.go"

	"github.com/luannevesbtc/TCStocksCrypto/app/crypto"
	"github.com/luannevesbtc/TCStocksCrypto/app/health"
	"github.com/luannevesbtc/TCStocksCrypto/store"
	auth "github.com/tradersclub/TCAuth/middleware/echo"
	"github.com/tradersclub/TCUtils/cache"
	"github.com/tradersclub/TCUtils/logger"
)

// Container modelo para exportação dos serviços instanciados
type Container struct {
	Health  health.App
	Crypto  crypto.App
	Session auth.Middleware
}

// Options struct de opções para a criação de uma instancia dos serviços
type Options struct {
	Stores  *store.Container
	Cache   cache.Cache
	Nats    *nats.Conn
	Session auth.Middleware

	StartedAt time.Time
	Version   string
}

// New cria uma nova instancia dos serviços
func New(opts Options) *Container {

	container := &Container{
		Health:  health.NewApp(opts.Stores, opts.Version, opts.StartedAt),
		Crypto:  crypto.NewApp(opts.Stores, opts.Nats, opts.Cache),
		Session: opts.Session,
	}

	logger.Info("Registered APP")

	return container

}
