package provider

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"github.com/go-redis/redis/v8"
	"server/pkg/app/model"
	"server/pkg/app/provider"
)

func NewTextProvider(textReadRepo model.TextReadRepository) provider.TextProvider {
	return &textProvider{
		textReadRepo: textReadRepo,
	}
}

type textProvider struct {
	textReadRepo model.TextReadRepository
}

type textSerializable struct {
	TextID     string  `json:"text_id"`
	Similarity int     `json:"similarity"`
	Value      string  `json:"value"`
	Rank       float64 `json:"rank"`
}

func (p *textProvider) GetTextID(text string) string {
	return hashText(text)
}

func (p *textProvider) GetByTextID(ctx context.Context, textID string) (provider.TextData, error) {
	text, err := p.textReadRepo.FindByID(ctx, model.TextID(textID))
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return provider.TextData{}, nil
		}
		return provider.TextData{}, err
	}

	return provider.TextData{
		Rank:       text.Rank(),
		Similarity: text.Similarity(),
		Value:      text.Value(),
	}, nil
}

func hashText(text string) string {
	hash := sha256.New()
	hash.Write([]byte(text))
	return hex.EncodeToString(hash.Sum(nil))
}
