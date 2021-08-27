package store

import (
	"github.com/jmoiron/sqlx"
	"github.com/tradersclub/TCTemplateBack/store/health"
	"github.com/tradersclub/TCTemplateBack/store/item"
	"github.com/tradersclub/TCUtils/logger"
)

// Container modelo para exportação dos repositórios instanciados
type Container struct {
	Health health.Store
	Item   item.Store
}

// Options struct de opções para a criação de uma instancia dos repositórios
type Options struct {
	Writer *sqlx.DB
	Reader *sqlx.DB
}

// New cria uma nova instancia dos repositórios
func New(opts Options) *Container {
	container := &Container{
		Health: health.NewStore(opts.Reader),
		Item:   item.NewStore(opts.Reader, opts.Writer),
	}

	logger.Info("Registered STORE")

	return container
}
