package server

import (
	"context"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/nats-io/nats.go"
	"github.com/tradersclub/TCUtils/cache"
	"github.com/tradersclub/TCUtils/logger"
	"github.com/tradersclub/TCUtils/natstan"
	"github.com/tradersclub/TCUtils/tcerr"
	"github.com/tradersclub/TCUtils/validator"

	"github.com/labstack/echo-contrib/prometheus"
	emiddleware "github.com/labstack/echo/v4/middleware"

	_ "github.com/go-sql-driver/mysql"
	"github.com/tradersclub/TCStocksCrypto/api"
	"github.com/tradersclub/TCStocksCrypto/app"
	pocConfig "github.com/tradersclub/TCStocksCrypto/config"
	"github.com/tradersclub/TCStocksCrypto/event"
	"github.com/tradersclub/TCStocksCrypto/model"
)

// Server is a interface to define contract to server up
type Server interface {
	Start()
	Stop()
	ReloadConnections()
}

type server struct {
	Echo       *echo.Echo
	Prometheus *prometheus.Prometheus
	Validator  *validator.Validator
	Nats       *nats.Conn

	Ctx context.Context

	App   *app.Container
	Cache cache.Cache
}

// New is instance the server
func New() Server {
	return &server{}
}

func (e *server) Start() {
	e.Echo = echo.New()
	e.Echo.Validator = validator.New()
	e.Echo.Debug = pocConfig.ConfigGlobal.ENV != "prod"
	e.Echo.HideBanner = true

	e.Echo.Use(emiddleware.Logger())
	e.Echo.Use(emiddleware.BodyLimit("2M"))
	e.Echo.Use(emiddleware.Recover())
	e.Echo.Use(emiddleware.RequestID())

	e.Prometheus = prometheus.NewPrometheus("TCStocksCrypto", nil)
	e.Prometheus.Use(e.Echo)

	e.StartNats()
	e.Cache = cache.NewMemcache(pocConfig.ConfigGlobal.Cache)
	e.StartApp()
	e.RegisterEvent()
	e.RegisterAPI()
	e.TreatErrorsHTTP()

	logger.Info("Start server PID: ", os.Getpid())
	if err := e.Echo.Start(pocConfig.ConfigGlobal.Server.Port); err != nil {
		logger.Error("cannot starting server ", err.Error())
	}
}

func (e *server) TreatErrorsHTTP() {
	// TODO: change to TCUtils
	e.Echo.HTTPErrorHandler = func(err error, c echo.Context) {
		if c.Response().Committed {
			return
		}

		if err := c.JSON(tcerr.GetHTTPCode(err), model.Response{Err: err}); err != nil {
			logger.ErrorContext(c.Request().Context(), err)
		}
	}
}

func (e *server) RegisterEvent() {
	event.Register(event.Options{
		Apps: e.App,
		Nats: e.Nats,
	})
}

func (e *server) StartApp() {
	e.App = app.New(app.Options{
		Cache: e.Cache,
		Nats:  e.Nats,
	})
}

func (e *server) RegisterAPI() {
	api.Register(api.Options{
		Group: e.Echo.Group(""),
		Apps:  e.App,
	})
}

func (e *server) StartNats() {
	e.Nats = natstan.New(natstan.Options{
		URL: pocConfig.ConfigGlobal.Nats.URL,
	})
}

func (e *server) Stop() {
	e.Nats.Close()
	if err := e.Echo.Close(); err != nil {
		logger.Error("cannot close echo ", err.Error())
	}
}

// ReloadConnections all connections like DB, Nats, ...
func (e *server) ReloadConnections() {
	e.Nats.Close()

	logger.Info("Close all connections...")
	e.StartNats()
}
