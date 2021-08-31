package api

import (
	"github.com/labstack/echo/v4"
	health "github.com/luannevesbtc/TCStocksCrypto/api/health"
	"github.com/luannevesbtc/TCStocksCrypto/app"
	"github.com/tradersclub/TCUtils/logger"
)

// Options struct de opções para a criação de uma instancia das rotas
type Options struct {
	Group *echo.Group
	Apps  *app.Container
}

// Register api instance
func Register(opts Options) {
	health.Register(opts.Group, opts.Apps)

	logger.Info("Registered API")
}
