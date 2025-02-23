package service

import (
	"context"
	"errors"
	"log"
	"server/pkg/app/model"
)

var ErrKeyAlreadyExists = errors.New("key already exists")

func NewValuatorService(repo model.TextRepository) ValuatorService {
	return &valuatorService{
		repo: repo,
	}
}

type ValuatorService interface {
	AddText(ctx context.Context, value string) (model.TextID, error)
}

type valuatorService struct {
	repo model.TextRepository
}

func (v *valuatorService) AddText(ctx context.Context, value string) (model.TextID, error) { // TODO без дубликатов и итерирования
	textID := v.repo.NextID(value)
	text := model.NewText(textID, value)

	err := v.repo.Store(ctx, text)
	if err != nil {
		if errors.Is(err, ErrKeyAlreadyExists) {
			return textID, err
		}
		log.Panic(err)
	}

	return textID, nil
}
