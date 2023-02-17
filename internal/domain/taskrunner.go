package domain

import "context"

// TaskRunnerService - provides access to task runner service.
type TaskRunnerService interface {
	RunTask(ctx context.Context, taskID TaskID) error
}
