package repo

import (
	"fmt"
	"github.com/go-redis/redis/v8"
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

func (sm *ShardManager) GetShard(region string) (*redis.Client, error) {
	shard, ok := sm.shards[region]
	if !ok {
		return nil, fmt.Errorf("no shard found for region %s", region)
	}
	return shard, nil
}

func (sm *ShardManager) GetShardByCountry(country string) (*redis.Client, error) {
	region := sm.regions[country]
	return sm.GetShard(region)
}
