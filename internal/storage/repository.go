package storage

import (
	"context"

	"github.com/redis/go-redis/v9"
)

type Repository interface {
	PFAddBatch(ctx context.Context, ips []string) error
}

type RedisClient struct {
	redisClient *redis.Client
}

func NewRedisClient(redisClient *redis.Client) *RedisClient {
	return &RedisClient{
		redisClient: redisClient,
	}
}

func (client *RedisClient) Client() *redis.Client {
	return client.redisClient
}
