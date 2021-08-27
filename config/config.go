package config

import (
	authClient "github.com/tradersclub/TCAuth/client"
	"github.com/tradersclub/TCUtils/cache"
)

// Docs is a struct to use in config
type Docs struct {
	Key string `json:"key"`
}

// Server is a struct to use in config
type Server struct {
	Port string `json:"port"`
}

// DBConn is a struct to use in Database
type DBConn struct {
	URL string `json:"url"`
}

// Database is a struct to use in config
type Database struct {
	Reader DBConn `json:"reader"`
	Writer DBConn `json:"Writer"`
}

// Nats is a struct to use in config
type Nats struct {
	URL string `json:"url"`
}

// TCAuth is a instance to valid Session
type TCAuth struct {
	Addr    string `json:"addr"`
	Port    string `json:"port"`
	Timeout string `json:"timeout"`
}

// Config is a struct to use in var ConfigGlobal
type Config struct {
	ENV      string            `json:"tc"`
	Docs     Docs              `json:"docs"`
	Server   Server            `json:"server"`
	Database Database          `json:"database"`
	Cache    cache.Options     `json:"cache"`
	Nats     Nats              `json:"nats"`
	Auth     authClient.Option `json:"auth"`
}

// ConfigGlobal is you use in all app
var ConfigGlobal *Config
