package repo

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v8"
	"log"
	"server/pkg/app/model"
)

func NewShardManager(
	mainClient *redis.Client,
	shards map[string]*redis.Client,
	regions map[string]string,
) *ShardManager {
	return &ShardManager{
		mainClient: mainClient,
		shards:     shards,
		regions:    regions,
	}
}

type ShardManager struct {
	mainClient *redis.Client
	shards     map[string]*redis.Client
	regions    map[string]string
}

func (sm *ShardManager) GetRepo(ctx context.Context, textID model.TextID) (model.TextRepository, error) {
	region, err := sm.mainClient.Get(ctx, fmt.Sprintf("text_region:%s", textID)).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, NotFoundRegion
		}
		return nil, err
	}
	shard, err := sm.getShard(region)
	if err != nil {
		return nil, err
	}

	log.Println()
	log.Printf("LOOKUP: %s, %s", textID, region)
	log.Println()

	return NewTextRepository(shard), nil
}

func (sm *ShardManager) getShard(region string) (*redis.Client, error) {
	shard, ok := sm.shards[region]
	if !ok {
		return nil, fmt.Errorf("no shard found for region %s", region)
	}

	return shard, nil
}
