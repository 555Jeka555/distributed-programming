package query

import (
	"context"

	"server/pkg/app/model"
)

type TextQueryService interface {
	GetTextByID(ctx context.Context, textID string, rankID string) (TextData, error)
}

func NewTextQueryService(repo model.TextReadRepository) TextQueryService {
	return &textQueryService{
		repo: repo,
	}
}

type TextData struct {
	Value string
	Rank  float64
}

type textQueryService struct {
	repo model.TextReadRepository
}

func (t *textQueryService) GetTextByID(ctx context.Context, textID string, rankID string) (TextData, error) {
	text, err := t.repo.FindByID(ctx, model.TextID(textID), model.RankID(rankID))
	if err != nil {
		return TextData{}, err
	}

	return TextData{
		Value: text.Value(),
		Rank:  text.Rank(),
	}, nil
}
