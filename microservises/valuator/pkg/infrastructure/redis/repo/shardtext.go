package repo

import (
	"context"
	"errors"
	"server/pkg/app/model"
)

var NotFoundRegion = errors.New("region not found")

func NewShardTextRepository(
	shardManager *ShardManager,
) model.TextRepository {
	return &shardTextRepository{
		shardManager: shardManager,
	}
}

type shardTextRepository struct {
	shardManager *ShardManager
}

func (t *shardTextRepository) GetTextID(text string) model.TextID {
	return model.TextID(hashText(text))
}

func (t *shardTextRepository) FindByID(ctx context.Context, textID model.TextID) (model.Text, error) {
	repo, err := t.shardManager.GetRepo(ctx, textID)
	if err != nil {
		return model.Text{}, err
	}

	return repo.FindByID(ctx, textID)
}
