package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/AssylzhanZharzhanov/axxonsoft-test-service/internal/domain"

	"github.com/go-redis/redis/v8"
)

type redisRepository struct {
	client *redis.Client
}

// NewRedisRepository - create a new redis repository
func NewRedisRepository(client *redis.Client) domain.TaskRedisRepository {
	return &redisRepository{
		client: client,
	}
}

func (r *redisRepository) Set(ctx context.Context, key string, value *domain.Task, seconds int) error {
	if len(key) == 0 {
		return fmt.Errorf("key required")
	}

	bytes, err := json.Marshal(value)
	if err != nil {
		return err
	}

	err = r.client.Set(ctx, key, bytes, time.Duration(seconds)).Err()
	if err != nil {
		return err
	}
	return nil
}

func (r *redisRepository) Get(ctx context.Context, key string) (*domain.Task, error) {
	value, err := r.client.Get(ctx, key).Bytes()
	if err != nil {
		return nil, err
	}

	var task *domain.Task
	err = json.Unmarshal(value, &task)
	if err != nil {
		return nil, err
	}

	return task, nil
}

func (r *redisRepository) Delete(ctx context.Context, key string) error {
	if err := r.client.Del(ctx, key).Err(); err != nil {
		return err
	}

	return nil
}
