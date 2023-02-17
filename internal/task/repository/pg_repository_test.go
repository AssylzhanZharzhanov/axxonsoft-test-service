package repository

import (
	"context"
	"testing"

	"github.com/AssylzhanZharzhanov/axxonsoft-test-service/internal/domain"
)

func TestRepository_Create(t *testing.T) {

	// Setup database.
	//
	runTestSetup(t, setupDatabase)
	defer purgeDatabase(testSQLDB)

	var (
		ctx = context.Background()
	)

	testRepository := NewRepository(database)

	// Prepare test data.
	//
	tests := []struct {
		name        string
		task        *domain.Task
		expectError bool
	}{
		{
			name: "Create task 1",
			task: &domain.Task{
				ID:       1,
				StatusID: 1,
				Method:   "GET",
				URL:      "test.com",
			},
			expectError: false,
		},
		{
			name: "Create task 2",
			task: &domain.Task{
				ID:       2,
				StatusID: 1,
				Method:   "GET",
				URL:      "test.com",
			},
			expectError: false,
		},
		{
			name: "Create task 3",
			task: &domain.Task{
				ID:       3,
				StatusID: 1,
				Method:   "GET",
				URL:      "test.com",
			},
			expectError: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			args := test.task
			result, err := testRepository.Create(ctx, args)
			if !test.expectError {
				if err != nil {
					t.Errorf("unexpected error %s", err)
				}
				if result == 0 {
					t.Error("task has not been created")
					return
				}
				if args.ID != result {
					t.Error("task has not been created")
				}
			} else {
				if err == nil {
					t.Error("expected error but got nothing")
				}
			}
		})
	}

}

func TestRepository_Get(t *testing.T) {

	// Setup database.
	//
	runTestSetup(t, setupDatabase)
	defer purgeDatabase(testSQLDB)

	var (
		ctx = context.Background()
	)

	testRepository := NewRepository(database)

	// Prepare test data.
	//
	tests := []struct {
		name        string
		task        *domain.Task
		expectError bool
	}{
		{
			name: "Create task 1",
			task: &domain.Task{
				ID:       1,
				StatusID: 1,
				Method:   "GET",
				URL:      "test.com",
			},
			expectError: false,
		},
		{
			name: "Create task 2",
			task: &domain.Task{
				ID:       2,
				StatusID: 1,
				Method:   "GET",
				URL:      "test.com",
			},
			expectError: false,
		},
		{
			name: "Create task 3",
			task: &domain.Task{
				ID:       3,
				StatusID: 1,
				Method:   "GET",
				URL:      "test.com",
			},
			expectError: false,
		},
	}

	// Create test items.
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			args := test.task
			_, err := testRepository.Create(ctx, args)
			if err != nil {
				t.Errorf("unexpected error %s", err)
			}
		})
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// create test items.
			task, err := testRepository.Get(ctx, test.task.ID)
			if err != nil {
				t.Errorf("Unexpected error %s", err)
			}
			if task.ID != test.task.ID {
				t.Errorf("user id is not same %s", err)
			}
		})
	}
}

func TestRepository_List(t *testing.T) {

	// Setup database.
	//
	runTestSetup(t, setupDatabase)
	defer purgeDatabase(testSQLDB)

	var (
		ctx = context.Background()
	)

	testRepository := NewRepository(database)

	// Prepare test data.
	//
	tests := []struct {
		name        string
		task        *domain.Task
		expectError bool
	}{
		{
			name: "Create task 1",
			task: &domain.Task{
				ID:       1,
				StatusID: 1,
				Method:   "GET",
				URL:      "test.com",
			},
			expectError: false,
		},
		{
			name: "Create task 2",
			task: &domain.Task{
				ID:       2,
				StatusID: 1,
				Method:   "GET",
				URL:      "test.com",
			},
			expectError: false,
		},
		{
			name: "Create task 3",
			task: &domain.Task{
				ID:       3,
				StatusID: 1,
				Method:   "GET",
				URL:      "test.com",
			},
			expectError: false,
		},
	}

	// Create test items.
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			args := test.task
			_, err := testRepository.Create(ctx, args)
			if err != nil {
				t.Errorf("unexpected error %s", err)
			}
		})
	}

	// Test pagination.
	list, total, err := testRepository.List(ctx, domain.TaskSearchCriteria{
		Page: domain.PageRequest{
			Offset: 0,
			Size:   3,
		},
	})
	if err != nil {
		t.Errorf("Unexpected error %s", err)
	}
	if len(list) != int(total) {
		t.Errorf("Pagination not working, expected items %d but got %v", 3, len(list))
	}
}
