package repo

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v8"
	"server/pkg/app/model"
	"server/pkg/app/service"
	"strconv"
)

const (
	textIDPattern = "TextID-%s"
	rankIDPattern = "RankID-%s"
)

func NewTextRepository(rdb *redis.Client) *TextRepository {
	return &TextRepository{
		rdb: rdb,
	}
}

type TextRepository struct {
	rdb *redis.Client
}

func (t *TextRepository) NextTextID(text string) model.TextID {
	return model.TextID(fmt.Sprintf(textIDPattern, hashText(text)))
}

func (t *TextRepository) NextRankID(text string) model.RankID {
	return model.RankID(fmt.Sprintf(rankIDPattern, hashText(text)))
}

func (t *TextRepository) Store(ctx context.Context, text model.Text) error {
	exists, err := t.keyExists(ctx, string(text.TextID()))
	if err != nil {
		return err
	}
	if exists {
		return service.ErrKeyAlreadyExists
	}

	exists, err = t.keyExists(ctx, string(text.RankID()))
	if err != nil {
		return err
	}
	if exists {
		return service.ErrKeyAlreadyExists
	}

	return t.storeText(ctx, text)
}

func (t *TextRepository) FindByID(ctx context.Context, textID model.TextID, rankID model.RankID) (model.Text, error) {
	value, err := t.rdb.Get(ctx, string(textID)).Result()
	if err != nil && !errors.Is(err, redis.Nil) {
		return model.Text{}, err
	}

	rankStr, err := t.rdb.Get(ctx, string(rankID)).Result()
	if err != nil && !errors.Is(err, redis.Nil) {
		return model.Text{}, err
	}
	rank, err := strconv.ParseFloat(rankStr, 64)
	if err != nil {
		return model.Text{}, fmt.Errorf("failed to parse rank: %w", err)
	}

	return model.LoadText(
		textID,
		rankID,
		value,
		rank,
	), nil
}

func (t *TextRepository) keyExists(ctx context.Context, key string) (bool, error) {
	_, err := t.rdb.Get(ctx, key).Result()
	if errors.Is(err, redis.Nil) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

func (t *TextRepository) storeText(ctx context.Context, text model.Text) error {
	_, err := t.rdb.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
		err := pipe.Set(ctx, string(text.TextID()), text.Value(), 0).Err()
		if err != nil {
			return err
		}

		err = pipe.Set(ctx, string(text.RankID()), text.Rank(), 0).Err()
		if err != nil {
			return err
		}

		return nil
	})
	return err
}

func hashText(text string) string {
	hash := sha256.New()
	hash.Write([]byte(text))
	return hex.EncodeToString(hash.Sum(nil))
}
