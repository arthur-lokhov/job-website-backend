// Package cache provides caching functionalities for the application.
package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

// NewRedisClient creates a new Redis client.
func NewRedisClient(redisURL string) (*redis.Client, error) {
	opt, err := redis.ParseURL(redisURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse redis URL: %w", err)
	}

	rdb := redis.NewClient(opt)

	// Ping the Redis server to check the connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err = rdb.Ping(ctx).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to redis: %w", err)
	}

	return rdb, nil
}
