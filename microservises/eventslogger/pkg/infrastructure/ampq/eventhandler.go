package ampq

import (
	"encoding/json"
	"fmt"
	"server/pkg/app/handler"

	"server/pkg/app/event"
)

type IntegrationEventHandler interface {
	Handle(body []byte) error
}

func NewIntegrationEventHandler(handler handler.Handler) IntegrationEventHandler {
	return &integrationEventHandler{
		handler: handler,
	}
}

type integrationEventHandler struct {
	handler handler.Handler
}

func (h *integrationEventHandler) Handle(body []byte) error {
	evt, err := parse(body)
	if err != nil {
		return err
	}

	return h.handler.Handle(evt)
}

func parse(body []byte) (event.Event, error) {
	var typeHolder struct {
		Type string `json:"type"`
	}

	if err := json.Unmarshal(body, &typeHolder); err != nil {
		return nil, fmt.Errorf("failed to unmarshal event type: %w", err)
	}

	switch typeHolder.Type {
	case "rankcalculator.rank_calculated":
		var msg RankCalculatedMessage
		if err := json.Unmarshal(body, &msg); err != nil {
			return nil, fmt.Errorf("failed to unmarshal rank calculated message: %w", err)
		}
		return event.RankCalculated{
			TextID: msg.TextID,
			Rank:   msg.Rank,
		}, nil
	case "valuator.similarity_calculated":
		var msg SimilarityCalculatedMessage
		if err := json.Unmarshal(body, &msg); err != nil {
			return nil, fmt.Errorf("failed to unmarshal similarity calculated message: %w", err)
		}
		return event.SimilarityCalculated{
			TextID:     msg.TextID,
			Similarity: msg.Similarity,
		}, nil

	default:
		return nil, fmt.Errorf("unknown event type: %s", typeHolder.Type)
	}
}

type RankCalculatedMessage struct {
	TextID string  `json:"text_id"`
	Rank   float64 `json:"rank"`
}

type SimilarityCalculatedMessage struct {
	TextID     string `json:"text_id"`
	Similarity int    `json:"similarity"`
}
