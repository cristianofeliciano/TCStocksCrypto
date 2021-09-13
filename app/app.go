package app

import (
	"time"

	"github.com/nats-io/nats.go"

	auth "github.com/tradersclub/TCAuth/middleware/echo"
	"github.com/tradersclub/TCStocksCrypto/app/crypto"
	"github.com/tradersclub/TCUtils/cache"
	"github.com/tradersclub/TCUtils/logger"
)

// Container modelo para exportação dos serviços instanciados
type Container struct {
	Crypto  crypto.App
	Session auth.Middleware
}

// Options struct de opções para a criação de uma instancia dos serviços
type Options struct {
	Cache   cache.Cache
	Nats    *nats.Conn
	Session auth.Middleware

	StartedAt time.Time
	Version   string
}

// New cria uma nova instancia dos serviços
func New(opts Options) *Container {

	container := &Container{
		Crypto:  crypto.NewApp(opts.Nats, opts.Cache),
		Session: opts.Session,
	}

	logger.Info("Registered APP")

	return container

}
