package handler

import (
	"context"
	"errors"
	"server/pkg/app/service"
)

type Handler interface {
	Handle(ctx context.Context, body string) error
}

func NewHandler(rankCalculator service.RankCalculatorService) Handler {
	return &handler{
		rankCalculator: rankCalculator,
	}
}

type handler struct {
	rankCalculator service.RankCalculatorService
}

func (h *handler) Handle(ctx context.Context, body string) error {
	err := h.rankCalculator.AddText(ctx, body)
	if err != nil && !errors.Is(err, service.ErrKeyAlreadyExists) {
		return err
	}

	return nil
}
