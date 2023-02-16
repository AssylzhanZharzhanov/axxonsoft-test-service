package service

import (
	"context"

	"github.com/AssylzhanZharzhanov/axxonsoft-test-service/internal/domain"

	"github.com/go-kit/log"
)

// Middleware describes a service (as opposed to endpoint) middleware.
type middleware func(taskService domain.TaskRunnerService) domain.TaskRunnerService

// loggingServiceMiddleware takes a logger as a dependency
// and returns a service Middleware.
func loggingServiceMiddleware(logger log.Logger) middleware {
	return func(next domain.TaskRunnerService) domain.TaskRunnerService {
		return loggingMiddleware{
			logger: logger,
			next:   next,
		}
	}
}

type loggingMiddleware struct {
	logger log.Logger
	next   domain.TaskRunnerService
}

func (mw loggingMiddleware) RunTask(ctx context.Context, taskID domain.TaskID) (err error) {
	defer func() {
		_ = mw.logger.Log("method", "RunTask",
			//domain.LogFieldTraceID, traceID,
			//domain.LogFieldSpanID, spanID,
			"taskID", taskID,
			"err", err,
		)
	}()
	return mw.next.RunTask(ctx, taskID)
}
