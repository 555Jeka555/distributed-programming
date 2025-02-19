package storage

import (
	"context"
	"errors"
	"github.com/go-redis/redis/v8"
	"log"
)

func NewKeyValue(rdb *redis.Client) *keyValue {
	return &keyValue{
		rdb: rdb,
	}
}

type keyValue struct {
	rdb *redis.Client
}

func (k *keyValue) Set(ctx context.Context, key string, text string) {
	err := k.rdb.Set(ctx, key, text, 0).Err()
	if err != nil {
		log.Fatalf("Could not add text to Redis: %v", err)
	}
}

func (k *keyValue) ListKey(ctx context.Context) ([]string, error) {
	var cursor uint64
	var keys []string
	for {
		result, nextCursor, err := k.rdb.Scan(ctx, cursor, "*", 0).Result()
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

func (k *keyValue) ListValue(ctx context.Context, keys []string) ([]string, error) {
	values, err := k.rdb.MGet(ctx, keys...).Result()
	if err != nil && !errors.Is(err, redis.Nil) {
		return nil, err
	}

	var result []string
	for _, value := range values {
		if value != nil {
			result = append(result, value.(string))
		}
	}
	return result, nil
}
