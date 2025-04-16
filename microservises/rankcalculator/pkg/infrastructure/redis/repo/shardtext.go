package repo

import (
	"context"
	"errors"
	"fmt"
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

func (t *shardTextRepository) Store(ctx context.Context, text model.Text) error {
	region, ok := RegionFromContext(ctx)
	if !ok {
		return errors.New(fmt.Sprintf("region not exists: %v", region))
	}

	err := t.shardManager.Store(ctx, string(text.TextID()), region, 0)
	if err != nil {
		return err
	}
	repo, err := t.shardManager.GetRepo(ctx, text.TextID())
	if err != nil {
		return err
	}

	return repo.Store(ctx, text)
}

func (t *shardTextRepository) Delete(ctx context.Context, textID model.TextID) error {
	repo, err := t.shardManager.GetRepo(ctx, textID)
	if err != nil {
		return err
	}

	return repo.Delete(ctx, textID)
}
