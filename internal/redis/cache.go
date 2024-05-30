package redis

import (
	"context"
	"time"

	rcl "github.com/redis/go-redis/v9"
)

type CacheConfig struct {
	Address string
	// Time-to-live for each cache entry before automatic deletion.
	TTL time.Duration
}

func NewCache(config CacheConfig) (*Client, error) {

	rdb := rcl.NewClient(&rcl.Options{
		Addr:             config.Address,
		Password:         "",
		DB:               0,
		DisableIndentity: true, // Disable set-info on connect
	})

	return &Client{
		redisClient: rdb,
	}, nil
}

type Client struct {
	config      CacheConfig
	redisClient *rcl.Client
}

func (cl *Client) Get(key string) ([]byte, error) {
	cmd := cl.redisClient.Get(context.Background(), key)
	if cmd.Err() != nil {
		return nil, cmd.Err()
	}
	return cmd.Bytes()
}
func (cl *Client) Set(key string, payload []byte) error {
	cmd := cl.redisClient.Set(context.Background(), key, payload, cl.config.TTL)
	return cmd.Err()
}
