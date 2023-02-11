package repository

import (
	"context"
	"log"
	"testing"

	"github.com/AssylzhanZharzhanov/axxonsoft-test-service/internal/domain"

	"github.com/alicebob/miniredis"
	"github.com/go-redis/redis/v8"
	"github.com/go-test/deep"
)

func SetupRedis() domain.TaskRedisRepository {
	miniRedis, err := miniredis.Run()
	if err != nil {
		log.Fatal(err)
	}

	client := redis.NewClient(&redis.Options{
		Addr: miniRedis.Addr(),
	})

	authRedisRepo := NewRedisRepository(client)
	return authRedisRepo
}

func TestRedisRepository_Set(t *testing.T) {

	var (
		ctx = context.Background()
	)

	redisRepoStub := SetupRedis()

	tests := []struct {
		name        string
		key         string
		value       *domain.Task
		expectError bool
	}{
		{
			name: "Success: task stored in cache",
			key:  "?/task/1",
			value: &domain.Task{
				Method: "GET",
				URL:    "test.com",
			},
			expectError: false,
		},
		{
			name: "Fail: key required",
			key:  "",
			value: &domain.Task{
				Method: "GET",
				URL:    "test.com",
			},
			expectError: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			key := test.key
			value := test.value

			err := redisRepoStub.Set(ctx, key, value, 60)
			if !test.expectError {
				if err != nil {
					t.Errorf("unexpected error: %s", err)
				}
			} else {
				if err == nil {
					t.Error("expected error but got nothing")
				}
			}
		})
	}
}

func TestRedisRepository_Get(t *testing.T) {

	var (
		ctx = context.Background()
	)

	redisRepoStub := SetupRedis()

	tests := []struct {
		name        string
		key         string
		value       *domain.Task
		expected    *domain.Task
		expectError bool
	}{
		{
			name: "Success: found in cache",
			key:  "?/task/1",
			value: &domain.Task{
				ID:     1,
				Method: "GET",
				URL:    "test.com",
			},
			expected: &domain.Task{
				ID:     1,
				Method: "GET",
				URL:    "test.com",
			},
			expectError: false,
		},
		{
			name:        "Fail: not found in cache",
			key:         "?task/2",
			value:       &domain.Task{},
			expected:    &domain.Task{},
			expectError: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			key := test.key
			value := test.value
			expected := test.expected

			err := redisRepoStub.Set(ctx, key, value, 60)
			if err != nil {
				t.Errorf("unexpected error: %s", err)
			}

			actual, err := redisRepoStub.Get(ctx, key)
			if !test.expectError {
				if err != nil {
					t.Errorf("unexpected error: %s", err)
				}
				if diff := deep.Equal(expected, actual); diff != nil {
					t.Error(diff)
				}
			} else {
				if err == nil {
					t.Error("expected error but got nothing")
				}
			}
		})
	}
}

func TestRedisRepository_Delete(t *testing.T) {

	var (
		ctx = context.Background()
	)

	redisRepoStub := SetupRedis()

	tests := []struct {
		name        string
		key         string
		value       *domain.Task
		expectError bool
	}{
		{
			name: "Success: found in cache",
			key:  "?/task/1",
			value: &domain.Task{
				ID:     1,
				Method: "GET",
				URL:    "test.com",
			},
			expectError: false,
		},
		{
			name:        "Fail: not found in cache",
			key:         "?task/2",
			value:       nil,
			expectError: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			key := test.key
			value := test.value

			err := redisRepoStub.Set(ctx, key, value, 60)
			if err != nil {
				t.Errorf("unexpected error: %s", err)
			}

			err = redisRepoStub.Delete(ctx, key)
			if err != nil {
				t.Errorf("unexpected error: %s", err)
			}

			actual, err := redisRepoStub.Get(ctx, key)
			if !test.expectError {
				if err.Error() != "redis: nil" {
					t.Errorf("unexpected error: %s", err)
				}
				if actual != nil {
					t.Errorf("value was not deleted")
				}
			} else {
				if err == nil {
					t.Error("expected error but got nothing")
				}
			}
		})
	}
}
