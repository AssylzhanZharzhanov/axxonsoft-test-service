package redis

import (
	"fmt"
	"time"

	"github.com/AssylzhanZharzhanov/axxonsoft-test-service/internal/domain"

	"github.com/go-redis/redis/v8"
)

// NewRedisClient - Returns new redis client
func NewRedisClient(cfg domain.AppConfig) (*redis.Client, error) {
	redisHost := fmt.Sprintf("%s:%s", cfg.RedisHost, cfg.RedisPort)

	client := redis.NewClient(&redis.Options{
		Addr:         redisHost,
		MinIdleConns: 200,
		PoolSize:     cfg.RedisPoolSize,
		PoolTimeout:  time.Duration(240) * time.Second,
		Password:     cfg.RedisPassword,
		DB:           cfg.RedisDB,
	})

	_, err := client.Ping(client.Context()).Result()
	if err != nil {
		return nil, err
	}

	return client, nil
}
