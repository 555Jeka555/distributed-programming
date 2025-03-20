package repo

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v8"
	"server/pkg/app/model"
	"strconv"
)

func NewTextRepository(rdb *redis.Client) *TextRepository {
	return &TextRepository{
		rdb: rdb,
	}
}

type TextRepository struct {
	rdb *redis.Client
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
