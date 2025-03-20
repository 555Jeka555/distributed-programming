package event

import (
	"context"
	"encoding/json"
	"errors"
	"server/pkg/app/service"
)

type Handler interface {
	Handle(ctx context.Context, body string) ([]byte, error)
}

func NewHandler(rankCalculator service.RankCalculatorService) Handler {
	return &handler{
		rankCalculator: rankCalculator,
	}
}

type handler struct {
	rankCalculator service.RankCalculatorService
}

type eventBody struct {
	TextID     string `json:"text_id"`
	RankID     string `json:"rank_id"`
	Similarity int    `json:"similarity"`
}

func (h *handler) Handle(ctx context.Context, body string) ([]byte, error) {
	textID, rankID, err := h.rankCalculator.AddText(ctx, body)
	similarity := 0
	if errors.Is(err, service.ErrKeyAlreadyExists) {
		similarity = 1
	}
	if err != nil && !errors.Is(err, service.ErrKeyAlreadyExists) {
		return nil, err
	}

	event := eventBody{
		TextID:     string(textID),
		RankID:     string(rankID),
		Similarity: similarity,
	}

	return json.Marshal(event)
}
