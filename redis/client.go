package redis

import (
	"context"
	"github.com/Hive-Gay/supreme-robot/config"
	"github.com/go-redis/redis/v8"
)

type Client struct {
	db  *redis.Client
	ctx context.Context
}

func NewClient(cfg *config.Config) (*Client, error) {
	client := Client{}

	// connect to redis
	client.db = redis.NewClient(&redis.Options{
		Addr:     cfg.RedisAddress,
		DB:       cfg.RedisDB,
		Password: cfg.RedisPassword,
	})

	// Create context
	client.ctx = context.Background()

	return &client, nil
}

func (c *Client) Client() (context.Context, *redis.Client) {
	return c.ctx, c.db
}
