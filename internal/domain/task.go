package domain

import (
	"context"
	"fmt"
)

const (
	TableName = "tasks"
)

// TaskID - id of the task
type TaskID int64

// StatusID - id of the status
type StatusID int64

// Header - represents header struct
type Header struct {
	Authentication string
	ContentType    string
}

// Task - represents task
type Task struct {
	ID             TaskID   `json:"ID"`
	Status         StatusID `json:"status"`
	HTTPStatusCode int      `json:"http_status_code"`
	ContentLength  int64    `json:"content_length"`
	Headers        Header   `json:"headers"`
	Method         string   `json:"method"`
	URL            string   `json:"URL"`
}

type TaskSearchCriteria struct {
}

// Validate - validates struct.
func (t Task) Validate() error {
	if len(t.Method) == 0 {
		return fmt.Errorf("invalid method")
	}
	if len(t.URL) == 0 {
		return fmt.Errorf("invlid url")
	}

	return nil
}

// TaskReadRepository - provides read access to a storage.
type TaskReadRepository interface {

	// Get - returns task from storage by id
	//
	Get(ctx context.Context, taskID TaskID) (*Task, error)

	// List - returns list of tasks from storage
	//
	List(ctx context.Context, criteria TaskSearchCriteria) ([]*Task, error)
}

// TaskRepository - provides access to a storage.
type TaskRepository interface {
	TaskReadRepository

	// Create - creates task in storage
	//
	Create(ctx context.Context, task *Task) (TaskID, error)
}

// TaskRedisRepository - provides access to a cache storage.
type TaskRedisRepository interface {

	// Set - stores task in cache for particular time
	//
	Set(ctx context.Context, key string, value *Task, seconds int) error

	// Get - returns task from cache by id
	//
	Get(ctx context.Context, key string) (*Task, error)

	// Delete - deletes from cache
	//
	Delete(ctx context.Context, key string) error
}

// TaskService - provides access to a business logic.
type TaskService interface {

	// CreateTask - creates task
	//
	CreateTask(ctx context.Context, task *Task) (TaskID, error)

	// GetTask - returns task by ID
	//
	GetTask(ctx context.Context, taskID TaskID) (*Task, error)

	// ListTasks - returns list of tasks
	//
	ListTasks(ctx context.Context, criteria TaskSearchCriteria) ([]*Task, error)
}
