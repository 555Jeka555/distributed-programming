package provider

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"github.com/go-redis/redis/v8"
	"server/pkg/app/provider"
	"server/pkg/infrastructure/keyvalue"
)

func NewTextProvider(rdb *redis.Client) provider.TextProvider {
	return &textProvider{
		storage: keyvalue.NewStorage[textSerializable](rdb),
	}
}

type textProvider struct {
	storage keyvalue.Storage[textSerializable]
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
	text, err := p.storage.Get(ctx, string(textID))
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return provider.TextData{}, nil
		}
		return provider.TextData{}, err
	}

	return provider.TextData{
		Rank:       text.Rank,
		Similarity: text.Similarity,
		Value:      text.Value,
	}, nil
}

func hashText(text string) string {
	hash := sha256.New()
	hash.Write([]byte(text))
	return hex.EncodeToString(hash.Sum(nil))
}
