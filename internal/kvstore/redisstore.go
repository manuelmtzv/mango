package kvstore

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisStore struct {
	client *redis.Client
}

func NewRedisStore(client *redis.Client) KVStorage {
	return &RedisStore{
		client: client,
	}
}

func (c *RedisStore) Set(ctx context.Context, key, val string, ttl time.Duration) error {
	return c.client.Set(ctx, key, val, ttl).Err()
}

func (c *RedisStore) Get(ctx context.Context, key string) (string, error) {
	return c.client.Get(ctx, key).Result()
}

func (c *RedisStore) Del(ctx context.Context, key string) error {
	return c.client.Del(ctx, key).Err()
}
