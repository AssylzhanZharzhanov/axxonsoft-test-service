package service

import (
	"context"
	"testing"

	"github.com/AssylzhanZharzhanov/axxonsoft-test-service/internal/domain"
	"github.com/AssylzhanZharzhanov/axxonsoft-test-service/internal/mocks"

	"github.com/go-test/deep"
	"github.com/golang/mock/gomock"
)

func TestService_CreateTask(t *testing.T) {

	var (
		ctx       = context.Background()
		validTask = &domain.Task{}
	)

	// Setup mocks.
	stubCtrl := gomock.NewController(t)
	defer stubCtrl.Finish()

	// Mocks repository.
	repositoryStub := mocks.NewMockTaskRepository(stubCtrl)
	eventServiceStub := mocks.NewMockEventService(stubCtrl)

	service := newBasicService(eventServiceStub, repositoryStub)

	eventServiceStub.EXPECT().
		RegisterEvent(ctx).
		Return().
		AnyTimes()

	repositoryStub.EXPECT().
		Create(ctx, validTask).
		Return(validTask, nil).
		AnyTimes()

	// Define tests.
	type arguments struct {
		task *domain.Task
	}

	type result struct {
		task *domain.Task
	}

	tests := []struct {
		name        string
		argument    arguments
		expected    result
		expectError bool
	}{
		{
			name: "Success: task created",
			argument: arguments{
				task: &domain.Task{},
			},
			expected: result{
				task: &domain.Task{},
			},
			expectError: false,
		},
		{
			name: "Fail: ",
			argument: arguments{
				task: &domain.Task{},
			},
			expected: result{
				task: &domain.Task{},
			},
			expectError: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			args := test.argument
			expected := test.expected

			taskID, err := service.CreateTask(ctx, args.task)
			if !test.expectError {
				if err != nil {
					t.Errorf("unexpected error: %s", err)
				}

				task, err := service.GetTask(ctx, taskID)
				if err != nil {
					t.Errorf("unexpected error: %s", err)
				}

				actual := result{
					task: task,
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
