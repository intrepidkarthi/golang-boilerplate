package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	"go-boilerplate/config"
	"go-boilerplate/internal/models"
	"time"
)

type RedisCache struct {
	client *redis.Client
}

func NewRedisCache(cfg *config.RedisConfig) (*RedisCache, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.Host, cfg.Port),
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	// Test connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to redis: %w", err)
	}

	return &RedisCache{client: client}, nil
}

func (c *RedisCache) SetMessage(ctx context.Context, message *models.Message) error {
	data, err := json.Marshal(message)
	if err != nil {
		return err
	}

	key := fmt.Sprintf("message:%s", message.ID.String())
	return c.client.Set(ctx, key, data, 24*time.Hour).Err()
}

func (c *RedisCache) GetMessage(ctx context.Context, id string) (*models.Message, error) {
	key := fmt.Sprintf("message:%s", id)
	data, err := c.client.Get(ctx, key).Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}
		return nil, err
	}

	var message models.Message
	if err := json.Unmarshal(data, &message); err != nil {
		return nil, err
	}

	return &message, nil
}
