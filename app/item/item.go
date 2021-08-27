package item

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/tradersclub/TCTemplateBack/model"
	"github.com/tradersclub/TCTemplateBack/store"
	"github.com/tradersclub/TCUtils/cache"
	"github.com/tradersclub/TCUtils/logger"
	"github.com/tradersclub/TCUtils/tcerr"
)

const TCSTARTKIT_GET_ITEM_BY_ID = "tcstartkit_get_item_by_id"

type getItemById struct {
	Err  error
	Id   string
	Data *model.Item
}

// App interface de item para implementação
type App interface {
	RequestItemById(ctx context.Context, id string) (*model.Item, error)
	GetItemById(ctx context.Context, id string) (*model.Item, error)
}

// NewApp cria uma nova instancia do serviço de exemplo item
func NewApp(stores *store.Container, nc *nats.Conn, cache cache.Cache) App {
	return &appImpl{
		stores: stores,
		nc:     nc,
		cache:  cache,
	}
}

type appImpl struct {
	stores    *store.Container
	startedAt time.Time
	version   string
	nc        *nats.Conn
	cache     cache.Cache
}

// RequestItemById - exemplo de request nats
func (s *appImpl) RequestItemById(ctx context.Context, id string) (*model.Item, error) {

	obj := new(getItemById)
	obj.Id = id
	jsonData, err := json.Marshal(obj)
	if err != nil {
		return nil, tcerr.NewError(http.StatusInternalServerError, "erro marshal item", nil)
	}
	msg, err := s.nc.Request(TCSTARTKIT_GET_ITEM_BY_ID, jsonData, time.Second*100)
	if err != nil {
		return nil, tcerr.NewError(http.StatusInternalServerError, "erro ao resgatar o item", nil)
	}

	var getItemByIdResponse *getItemById
	errResponse := json.Unmarshal(msg.Data, &getItemByIdResponse)

	if errResponse != nil {
		return nil, tcerr.NewError(http.StatusInternalServerError, "erro ao realizar o parser do item", nil)
	}

	return getItemByIdResponse.Data, nil
}

// GetItemById - exemplo de recuperação de dados do cache e store
func (s *appImpl) GetItemById(ctx context.Context, id string) (*model.Item, error) {
	item := new(model.Item)

	// exemplo de consulta em cache caso seja necessário, não esquece de validar junto ao seu líder qual memcached usar
	if err := s.cache.Get(ctx, id, item); err != nil {
		logger.ErrorContext(ctx, "app.item.get_item_by_id", "não encontrei o cache com id: "+id, err.Error())
	}

	// exemplo de consulta em store
	result := <-s.stores.Item.GetItemById(ctx, id)
	if result.Error != nil {
		logger.ErrorContext(ctx, "app.item.get_item_by_id", result.Error.Error())
		return nil, result.Error
	}

	data, err := model.ToItem(result.Data)
	if err != nil {
		logger.ErrorContext(ctx, "app.item.get_item_by_id", err.Error())

		return nil, err
	}

	return data, nil
}
