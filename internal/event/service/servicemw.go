package service

import (
	"context"

	"github.com/AssylzhanZharzhanov/axxonsoft-test-service/internal/domain"

	"github.com/go-kit/log"
)

// Middleware describes a service (as opposed to endpoint) middleware.
type middleware func(taskService domain.EventService) domain.EventService

// loggingServiceMiddleware takes a logger as a dependency
// and returns a service Middleware.
func loggingServiceMiddleware(logger log.Logger) middleware {
	return func(next domain.EventService) domain.EventService {
		return loggingMiddleware{
			logger: logger,
			next:   next,
		}
	}
}

type loggingMiddleware struct {
	logger log.Logger
	next   domain.EventService
}

func (mw loggingMiddleware) RegisterEvent(ctx context.Context) (err error) {
	defer func() {
		_ = mw.logger.Log("method", "RegisterEvent",
			//domain.LogFieldTraceID, traceID,
			//domain.LogFieldSpanID, spanID,
			//"task", task,
			"err", err,
		)
	}()
	return mw.next.RegisterEvent(ctx)
}
