package cache

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type Cache struct {
	client *redis.Client
}

func New(addr, password string, db int) (*Cache, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := rdb.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to redis: %w", err)
	}

	log.Println("Connected to Redis")
	return &Cache{client: rdb}, nil
}

func (c *Cache) Close() {
	if c.client != nil {
		c.client.Close()
		log.Println("Redis connection closed")
	}
}

func (c *Cache) CreateSession(ctx context.Context, userID uuid.UUID, ttl time.Duration) (string, error) {
	sessionID := uuid.New().String()
	if err := c.client.Set(ctx, sessionID, userID.String(), ttl).Err(); err != nil {
		return "", err
	}
	return sessionID, nil
}

func (c *Cache) GetSession(ctx context.Context, sessionID string) (uuid.UUID, error) {
	userIDStr, err := c.client.Get(ctx, sessionID).Result()
	if err != nil {
		return uuid.Nil, err
	}
	return uuid.Parse(userIDStr)
}

func (c *Cache) DeleteSession(ctx context.Context, sessionID string) error {
	return c.client.Del(ctx, sessionID).Err()
}
