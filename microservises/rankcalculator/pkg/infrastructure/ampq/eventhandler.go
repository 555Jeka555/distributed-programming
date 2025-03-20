package ampq

import (
	"context"
	"encoding/json"

	"server/pkg/app/event"
)

type IntegrationEventHandler interface {
	Handle(ctx context.Context, body []byte) ([]byte, error)
}

func NewIntegrationEventHandler(handler event.Handler) IntegrationEventHandler {
	return &integrationEventHandler{
		handler: handler,
	}
}

type integrationEventHandler struct {
	handler event.Handler
}

type eventBody struct {
	Text string `json:"text"`
}

func (h *integrationEventHandler) Handle(ctx context.Context, body []byte) ([]byte, error) {
	var evt eventBody

	err := json.Unmarshal(body, &evt)
	if err != nil {
		return nil, err
	}

	return h.handler.Handle(ctx, evt.Text)
}
