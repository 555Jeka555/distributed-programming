package repo

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"github.com/go-redis/redis/v8"
	"log"
	"server/pkg/app/model"
	"server/pkg/app/service"
)

func NewTextRepository(rdb *redis.Client) *TextRepository {
	return &TextRepository{
		rdb: rdb,
	}
}

type TextRepository struct {
	rdb *redis.Client
}

func (t *TextRepository) NextID(text string) model.TextID {
	return model.TextID(hashText(text))
}

func (t *TextRepository) Store(ctx context.Context, text model.Text) error {
	_, err := t.rdb.Get(ctx, string(text.ID())).Result()
	if !errors.Is(err, redis.Nil) {
		return service.ErrKeyAlreadyExists
	}
	if err != nil && !errors.Is(err, redis.Nil) {
		return err
	}

	err = t.rdb.Set(ctx, string(text.ID()), text.Value(), 0).Err()
	if err != nil && !errors.Is(err, redis.Nil) {
		return err
	}

	return nil
}

func (t *TextRepository) FindByID(ctx context.Context, textID model.TextID) (model.Text, error) {
	value, err := t.rdb.Get(ctx, string(textID)).Result()
	if err != nil && !errors.Is(err, redis.Nil) {
		return model.Text{}, err
	}

	return model.LoadText(
		textID,
		value,
	), nil
}

func (t *TextRepository) ListAll(ctx context.Context) ([]model.Text, error) {
	keys, err := t.listKey(ctx)
	if err != nil {
		log.Panic(err)
	}

	texts := make([]model.Text, 0, len(keys))
	for _, key := range keys {
		value, err1 := t.rdb.Get(ctx, key).Result()
		if err1 != nil && !errors.Is(err1, redis.Nil) {
			log.Panic(err1)
			return nil, err1
		}

		text := model.LoadText(model.TextID(key), value)
		texts = append(texts, text)
	}

	return texts, nil
}

func (t *TextRepository) listKey(ctx context.Context) ([]string, error) {
	var cursor uint64
	var keys []string
	for {
		result, nextCursor, err := t.rdb.Scan(ctx, cursor, "*", 0).Result()
		if err != nil && !errors.Is(err, redis.Nil) {
			return nil, err
		}
		keys = append(keys, result...)
		cursor = nextCursor
		if cursor == 0 {
			break
		}
	}
	return keys, nil
}

func hashText(text string) string {
	hash := sha256.New()
	hash.Write([]byte(text))
	return hex.EncodeToString(hash.Sum(nil))
}
