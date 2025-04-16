package repo

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v8"
	"log"
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
	repo, err := t.getRepo(ctx, textID)
	if err != nil {
		return model.Text{}, err
	}

	return repo.FindByID(ctx, textID)
}

func (t *shardTextRepository) getRepo(ctx context.Context, textID model.TextID) (model.TextRepository, error) {
	region, err := t.shardManager.mainClient.Get(ctx, fmt.Sprintf("text_region:%s", textID)).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, NotFoundRegion
		}
		return nil, err
	}
	shard, err := t.shardManager.GetShard(region)
	if err != nil {
		return nil, err
	}
	log.Println()
	log.Printf("LOOKUP: %s, %s", textID, region)
	log.Println()

	return NewTextRepository(shard), nil
}
