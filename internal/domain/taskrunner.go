package domain

import "context"

type TaskRunnerService interface {
	RunTask(ctx context.Context, taskID TaskID) error
}
