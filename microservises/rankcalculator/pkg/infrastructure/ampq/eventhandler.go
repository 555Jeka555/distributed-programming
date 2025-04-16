package ampq

import (
	"context"
	"encoding/json"
	"server/pkg/app/handler"
	"server/pkg/infrastructure/redis/repo"
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
	Text   string `json:"text"`
	Region string `json:"region"`
}

func (h *integrationEventHandler) Handle(ctx context.Context, body []byte) error {
	var evt eventBody

	err := json.Unmarshal(body, &evt)
	if err != nil {
		return err
	}

	return h.handler.Handle(context.WithValue(ctx, repo.RegionKey{}, evt.Region), evt.Text)
}
