package handler

import (
	"context"
	"errors"
	"math/rand"
	"server/pkg/app/provider"
	"time"

	"server/pkg/app/service"
)

type Handler interface {
	Handle(ctx context.Context, body string) error
}

func NewHandler(
	rankCalculator service.RankCalculatorService,
	textProvider provider.TextProvider,
	centrifugoClient service.CentrifugoClient,
) Handler {
	return &handler{
		rankCalculator:   rankCalculator,
		textProvider:     textProvider,
		centrifugoClient: centrifugoClient,
	}
}

type handler struct {
	rankCalculator   service.RankCalculatorService
	textProvider     provider.TextProvider
	centrifugoClient service.CentrifugoClient
}

func (h *handler) Handle(ctx context.Context, body string) error {
	delay := time.Duration(rand.Intn(2)+3) * time.Second
	time.Sleep(delay)

	err := h.rankCalculator.AddText(ctx, body)
	if err != nil && !errors.Is(err, service.ErrKeyAlreadyExists) {
		return err
	}

	textID := h.textProvider.GetTextID(body)
	textData, err := h.textProvider.GetByTextID(ctx, textID)
	if err != nil {
		return err
	}

	channel := "results"
	return h.centrifugoClient.Publish(channel, map[string]interface{}{
		"textID":     textID,
		"textValue":  textData.Value,
		"similarity": textData.Similarity,
		"rank":       textData.Rank,
	})
}
