package server

import (
	"context"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/nats-io/nats.go"
	"github.com/tradersclub/TCTemplateBack/api/swagger"
	"github.com/tradersclub/TCTemplateBack/store"
	"github.com/tradersclub/TCUtils/cache"
	"github.com/tradersclub/TCUtils/logger"
	"github.com/tradersclub/TCUtils/natstan"
	"github.com/tradersclub/TCUtils/tcerr"
	"github.com/tradersclub/TCUtils/validator"

	"github.com/labstack/echo-contrib/prometheus"
	emiddleware "github.com/labstack/echo/v4/middleware"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/tradersclub/TCTemplateBack/api"
	"github.com/tradersclub/TCTemplateBack/app"
	pocConfig "github.com/tradersclub/TCTemplateBack/config"
	"github.com/tradersclub/TCTemplateBack/event"
	"github.com/tradersclub/TCTemplateBack/model"

	auth "github.com/tradersclub/TCAuth/middleware/echo"
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
	Session    auth.Middleware

	DBReader     *sqlx.DB
	StopDBReader context.CancelFunc
	DBWriter     *sqlx.DB
	StopDBWriter context.CancelFunc
	Store        *store.Container

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

	e.Prometheus = prometheus.NewPrometheus("TCTemplateBack", nil)
	e.Prometheus.Use(e.Echo)

	e.StartDB()
	e.StartNats()
	e.Cache = cache.NewMemcache(pocConfig.ConfigGlobal.Cache)
	e.StartTCAuthClient()
	e.StartApp()
	e.RegisterEvent()
	e.RegisterAPI()
	e.TreatErrorsHTTP()

	if e.Echo.Debug {
		swagger.Register(swagger.Options{
			Port:      pocConfig.ConfigGlobal.Server.Port,
			Group:     e.Echo.Group("/swagger"),
			AccessKey: pocConfig.ConfigGlobal.Docs.Key,
		})
	}

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

// StartTCAuthClient only instance CientSession with middleware
func (e *server) StartTCAuthClient() {
	e.Session = auth.NewMiddle(pocConfig.ConfigGlobal.Auth)
}

func (e *server) CloseTCAuth() {
	if e.Session != nil {
		if err := e.Session.Close(); err != nil {
			logger.Error("Cannot close gRPC from TCAuth ", err.Error())
		} else {
			logger.Info("Closed gRPC from TCAuth")
		}
	}
}

func (e *server) ConnectTCAuthClient() {
	e.Session.Connect(pocConfig.ConfigGlobal.Auth)
	logger.Info("Connected TCAuth via gRPC")
}

func (e *server) RegisterEvent() {
	event.Register(event.Options{
		Apps: e.App,
		Nats: e.Nats,
	})
}

func (e *server) StartApp() {
	e.App = app.New(app.Options{
		Stores:  e.Store,
		Cache:   e.Cache,
		Nats:    e.Nats,
		Session: e.Session,
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

func (e *server) startDBReader() {
	var err error
	e.Ctx, e.StopDBReader = context.WithCancel(context.Background())
	e.DBReader, err = sqlx.ConnectContext(e.Ctx, "mysql", pocConfig.ConfigGlobal.Database.Reader.URL)
	if err != nil {
		logger.ErrorContext(e.Ctx, "Not connect DB Reader: ", err.Error())
		return
	}
	logger.Info("DB Reader was connected")
}

func (e *server) startDBWriter() {
	var err error
	e.Ctx, e.StopDBWriter = context.WithCancel(context.Background())
	e.DBWriter, err = sqlx.ConnectContext(e.Ctx, "mysql", pocConfig.ConfigGlobal.Database.Writer.URL)
	if err != nil {
		logger.ErrorContext(e.Ctx, "Not connect DB Writer: ", err.Error())
		return
	}
	logger.Info("DB Writer was connected")
}

func (e *server) StartDB() {
	e.startDBReader()
	e.startDBWriter()
	e.Store = store.New(store.Options{
		Writer: e.DBWriter,
		Reader: e.DBReader,
	})
}

func (e *server) Stop() {
	e.StopDB()
	e.Nats.Close()
	e.CloseTCAuth()
	if err := e.Echo.Close(); err != nil {
		logger.Error("cannot close echo ", err.Error())
	}
}

func (e *server) StopDB() {
	if e.DBReader != nil {
		if err := e.DBReader.Close(); err != nil {
			logger.Error("Cannot close dabatabse Reader: ", err.Error())
		} else {
			e.StopDBReader()
			logger.Info("Closed DBReader")
		}
	}
	if e.DBWriter != nil {
		if err := e.DBWriter.Close(); err != nil {
			logger.Error("Cannot close dabatabse Writer: ", err.Error())
		} else {
			e.StopDBWriter()
			logger.Info("Closed DBWriter")
		}
	}
}

// ReloadConnections all connections like DB, Nats, ...
func (e *server) ReloadConnections() {
	e.StopDB()
	e.Nats.Close()
	e.CloseTCAuth()

	logger.Info("Close all connections...")
	e.StartDB()
	e.StartNats()
	e.ConnectTCAuthClient()

}
