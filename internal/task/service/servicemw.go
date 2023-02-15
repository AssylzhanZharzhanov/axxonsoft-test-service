package service

import (
	"context"

	"github.com/AssylzhanZharzhanov/axxonsoft-test-service/internal/domain"

	"github.com/go-kit/log"
)

// Middleware describes a service (as opposed to endpoint) middleware.
type middleware func(taskService domain.TaskService) domain.TaskService

// loggingServiceMiddleware takes a logger as a dependency
// and returns a service Middleware.
func loggingServiceMiddleware(logger log.Logger) middleware {
	return func(next domain.TaskService) domain.TaskService {
		return loggingMiddleware{
			logger: logger,
			next:   next,
		}
	}
}

type loggingMiddleware struct {
	logger log.Logger
	next   domain.TaskService
}

func (mw loggingMiddleware) CreateTask(ctx context.Context, task *domain.Task) (result domain.TaskID, err error) {
	defer func() {
		_ = mw.logger.Log("method", "CreateTask",
			//domain.LogFieldTraceID, traceID,
			//domain.LogFieldSpanID, spanID,
			"task", task,
			"err", err,
			"result", result)
	}()
	return mw.next.CreateTask(ctx, task)
}

func (mw loggingMiddleware) GetTask(ctx context.Context, taskID domain.TaskID) (result *domain.Task, err error) {
	defer func() {
		_ = mw.logger.Log("method", "GetTask",
			//domain.LogFieldTraceID, traceID,
			//domain.LogFieldSpanID, spanID,
			"taskID", taskID,
			"err", err,
			"result", result)
	}()
	return mw.next.GetTask(ctx, taskID)
}

func (mw loggingMiddleware) ListTasks(ctx context.Context, criteria domain.TaskSearchCriteria) (result []*domain.Task, total domain.Total, err error) {
	defer func() {
		_ = mw.logger.Log("method", "ListTasks",
			//domain.LogFieldTraceID, traceID,
			//domain.LogFieldSpanID, spanID,
			"criteria", criteria,
			"err", err,
			"result", result,
			"total", total,
		)
	}()
	return mw.next.ListTasks(ctx, criteria)
}
