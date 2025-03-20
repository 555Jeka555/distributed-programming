package repo

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"github.com/go-redis/redis/v8"
	"server/pkg/app/model"
	"server/pkg/infrastructure/keyvalue"
)

func NewTextRepository(rdb *redis.Client) model.TextRepository {
	return &textRepository{
		storage: keyvalue.NewStorage[textSerializable](rdb),
	}
}

type textSerializable struct {
	TextID     string  `json:"text_id"`
	Similarity int     `json:"similarity"`
	Value      string  `json:"value"`
	Rank       float64 `json:"rank"`
}

type textRepository struct {
	storage keyvalue.Storage[textSerializable]
}

func (t *textRepository) NextTextID(text string) model.TextID {
	return model.TextID(hashText(text))
}

func (t *textRepository) FindByID(ctx context.Context, textID model.TextID) (model.Text, error) {
	text, err := t.storage.Get(ctx, string(textID))
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return model.Text{}, nil
		}
		return model.Text{}, err
	}

	return model.LoadText(
		model.TextID(text.TextID),
		text.Similarity,
		text.Value,
		text.Rank,
	), nil
}

func hashText(text string) string {
	hash := sha256.New()
	hash.Write([]byte(text))
	return hex.EncodeToString(hash.Sum(nil))
}
