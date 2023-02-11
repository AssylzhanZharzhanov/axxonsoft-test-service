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
		validTask = &domain.Task{
			ID:             1,
			Status:         1,
			HTTPStatusCode: 0,
			ContentLength:  0,
			Method:         "GET",
			URL:            "test.com",
			Headers: domain.Header{
				Authentication: "Bearer 123",
			},
		}
	)

	// Setup mocks.
	stubCtrl := gomock.NewController(t)
	defer stubCtrl.Finish()

	// Mocks repository.
	repositoryStub := mocks.NewMockTaskRepository(stubCtrl)
	redisStub := mocks.NewMockTaskRedisRepository(stubCtrl)
	eventServiceStub := mocks.NewMockEventService(stubCtrl)

	service := newBasicService(eventServiceStub, repositoryStub, redisStub)

	//eventServiceStub.EXPECT().
	//	RegisterEvent(ctx).
	//	Return().
	//	AnyTimes()

	repositoryStub.EXPECT().
		Create(ctx, validTask).
		Return(validTask.ID, nil).
		AnyTimes()

	// Define tests.
	type arguments struct {
		task *domain.Task
	}

	type result struct {
		taskID domain.TaskID
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
				task: validTask,
			},
			expected: result{
				taskID: validTask.ID,
			},
			expectError: false,
		},
		{
			name: "Fail: method required",
			argument: arguments{
				task: &domain.Task{
					Method: "",
					URL:    "test.com",
					Headers: domain.Header{
						Authentication: "Bearer 123",
					},
				},
			},
			expected: result{
				taskID: domain.TaskID(0),
			},
			expectError: true,
		},
		{
			name: "Fail: url required",
			argument: arguments{
				task: &domain.Task{
					Method: "GET",
					URL:    "",
					Headers: domain.Header{
						Authentication: "Bearer 123",
					},
				},
			},
			expected: result{
				taskID: domain.TaskID(0),
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

				actual := result{
					taskID: taskID,
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

func TestService_GetTask(t *testing.T) {
	var (
		ctx        = context.Background()
		cachedKey  = "?task/1"
		validTask1 = &domain.Task{
			ID:             1,
			Status:         3,
			HTTPStatusCode: 200,
			ContentLength:  200,
			Method:         "GET",
			URL:            "test.com",
		}

		notCachedKey = "?task/2"
		validTask2   = &domain.Task{
			ID:             2,
			Status:         3,
			HTTPStatusCode: 200,
			ContentLength:  200,
			Method:         "POST",
			URL:            "test.com",
		}
	)

	// Setup mocks.
	stubCtrl := gomock.NewController(t)
	defer stubCtrl.Finish()

	// Mocks repository.
	repositoryStub := mocks.NewMockTaskRepository(stubCtrl)
	redisStub := mocks.NewMockTaskRedisRepository(stubCtrl)
	eventServiceStub := mocks.NewMockEventService(stubCtrl)

	service := newBasicService(eventServiceStub, repositoryStub, redisStub)

	redisStub.EXPECT().
		Get(ctx, cachedKey).
		Return(validTask1, nil).
		AnyTimes()

	repositoryStub.EXPECT().
		Get(ctx, validTask1.ID).
		Return(validTask1, nil).
		AnyTimes()

	redisStub.EXPECT().
		Get(ctx, notCachedKey).
		Return(nil, nil).
		AnyTimes()

	repositoryStub.EXPECT().
		Get(ctx, validTask2.ID).
		Return(validTask2, nil).
		AnyTimes()

	redisStub.EXPECT().
		Set(ctx, notCachedKey, validTask2, 3600).
		Return(nil).
		AnyTimes()

	// Define tests.
	type arguments struct {
		taskID domain.TaskID
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
			name: "Success: task was found",
			argument: arguments{
				taskID: domain.TaskID(1),
			},
			expected: result{
				task: &domain.Task{},
			},
			expectError: false,
		},
		{
			name: "Fail: invalid task id",
			argument: arguments{
				taskID: domain.TaskID(0),
			},
			expected: result{
				task: &domain.Task{},
			},
			expectError: true,
		},
		{
			name: "Success: get data from storage",
			argument: arguments{
				taskID: domain.TaskID(2),
			},
			expected: result{
				task: validTask2,
			},
			expectError: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			args := test.argument
			expected := test.expected

			task, err := service.GetTask(ctx, args.taskID)
			if !test.expectError {
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
