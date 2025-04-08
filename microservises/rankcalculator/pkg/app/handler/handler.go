package handler

import (
	"context"
	"errors"
	"fmt"
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
	delay := time.Duration(3) * time.Second
	time.Sleep(delay)

	err := h.rankCalculator.AddText(ctx, body)
	if err != nil && !errors.Is(err, service.ErrKeyAlreadyExists) {
		return err
	}

	textID := h.textProvider.GetTextID(body)
	fmt.Println("textID", textID)
	fmt.Println("textID", textID)
	fmt.Println("textID", textID)

	channel := "results"
	err = h.centrifugoClient.Publish(channel, map[string]interface{}{
		"textID":     textID,
		"similarity": 12,
		"rank":       23,
	})
	if err != nil {
		return err
	}

	return nil
}
