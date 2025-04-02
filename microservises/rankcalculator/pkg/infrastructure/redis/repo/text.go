package repo

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
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

func (t *textRepository) GetTextID(text string) model.TextID {
	return model.TextID(hashText(text))
}

func (t *textRepository) Store(ctx context.Context, text model.Text) error {
	return t.storage.Set(ctx, string(text.TextID()), textSerializable{
		TextID:     string(text.TextID()),
		Similarity: text.Similarity(),
		Rank:       text.Rank(),
		Value:      text.Value(),
	}, 0)
}

func (t *textRepository) Delete(ctx context.Context, textID model.TextID) error {
	return t.storage.Delete(ctx, string(textID))
}

func hashText(text string) string {
	hash := sha256.New()
	hash.Write([]byte(text))
	return hex.EncodeToString(hash.Sum(nil))
}
