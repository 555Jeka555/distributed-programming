package repo

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

type RegionKey struct{}

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

func (sm *ShardManager) Store(ctx context.Context, textID, region string, expiration int) error {
	_, ok := sm.shards[region]
	if !ok {
		return fmt.Errorf("no shard found for region %s", region)
	}

	return sm.mainClient.Set(ctx, fmt.Sprintf("text_region:%s", textID), region, time.Duration(expiration)).Err()
}

func (sm *ShardManager) GetShard(region string) (*redis.Client, error) {
	shard, ok := sm.shards[region]
	if !ok {
		return nil, fmt.Errorf("no shard found for region %s", region)
	}

	return shard, nil
}

func RegionFromContext(ctx context.Context) (string, bool) {
	region, ok := ctx.Value(RegionKey{}).(string)

	return region, ok
}
