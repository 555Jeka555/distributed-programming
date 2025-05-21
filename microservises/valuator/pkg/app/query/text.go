package query

import (
	"context"

	"server/pkg/app/model"
)

type TextQueryService interface {
	GetTextByID(ctx context.Context, textID string) (TextData, error)
}

func NewTextQueryService(repo model.TextReadRepository) TextQueryService {
	return &textQueryService{
		repo: repo,
	}
}

type TextData struct {
	Similarity int
	Value      string
	Login      string
	Rank       float64
}

type textQueryService struct {
	repo model.TextReadRepository
}

func (t *textQueryService) GetTextByID(ctx context.Context, textID string) (TextData, error) {
	text, err := t.repo.FindByID(ctx, model.TextID(textID))
	if err != nil {
		return TextData{}, err
	}

	return TextData{
		Similarity: text.Similarity(),
		Value:      text.Value(),
		Login:      text.Login(),
		Rank:       text.Rank(),
	}, nil
}
