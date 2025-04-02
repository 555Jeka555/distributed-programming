package ampq

import (
	"context"
	"encoding/json"
	"server/pkg/app/handler"
)

type IntegrationEventHandler interface {
	Handle(ctx context.Context, body []byte) error
}

func NewIntegrationEventHandler(handler handler.Handler) IntegrationEventHandler {
	return &integrationEventHandler{
		handler: handler,
	}
}

type integrationEventHandler struct {
	handler handler.Handler
}

type eventBody struct {
	Text string `json:"text"`
}

func (h *integrationEventHandler) Handle(ctx context.Context, body []byte) error {
	var evt eventBody

	err := json.Unmarshal(body, &evt)
	if err != nil {
		return err
	}

	return h.handler.Handle(ctx, evt.Text)
}
