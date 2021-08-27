package item

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/nats-io/nats.go"

	"github.com/tradersclub/TCTemplateBack/app"
	"github.com/tradersclub/TCTemplateBack/model"
	"github.com/tradersclub/TCUtils/logger"
	"github.com/tradersclub/TCUtils/tcerr"
)

const TCSTARTKIT_GET_ITEM_BY_ID = "tcstartkit_get_item_by_id"

type getItemById struct {
	Err  error
	Id   string
	Data *model.Item
}

// Register group health check
func Register(apps *app.Container, conn *nats.Conn) {
	e := &event{
		apps: apps,
		nc:   conn,
	}

	e.nc.Subscribe(TCSTARTKIT_GET_ITEM_BY_ID, e.getItemById)
}

type event struct {
	apps *app.Container
	nc   *nats.Conn
}

func (e *event) getItemById(msg *nats.Msg) {
	ctx := context.Background()

	//recupera o objeto através do metódo criado no TCNatsModel
	var getItemByIdResponse *getItemById
	errResponse := json.Unmarshal(msg.Data, &getItemByIdResponse)

	// valida erro de parser
	if errResponse != nil {
		logger.ErrorContext(ctx, "event.session.get_session_by_id", errResponse.Error())
		obj := new(getItemById)
		obj.Err = tcerr.NewError(500, "event.session.get_session_by_id", errResponse.Error())
		b, _ := json.Marshal(obj)
		msg.Respond(b)
		return
	}

	// pega o item
	session, err := e.apps.Item.GetItemById(ctx, getItemByIdResponse.Id)
	if err != nil {
		logger.ErrorContext(ctx, "event.session.get_session_by_id", err.Error())
		getItemByIdResponse.Err = tcerr.NewError(http.StatusInternalServerError, "event.session.get_session_by_id", err.Error())
	} else {
		getItemByIdResponse.Data = session
	}

	b, _ := json.Marshal(getItemByIdResponse)
	msg.Respond(b)
}
