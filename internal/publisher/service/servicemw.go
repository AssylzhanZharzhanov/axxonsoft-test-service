package service

import (
	"context"

	"github.com/AssylzhanZharzhanov/axxonsoft-test-service/internal/domain"

	"github.com/go-kit/log"
)

// Middleware describes a service (as opposed to endpoint) middleware.
type middleware func(taskService domain.Publisher) domain.Publisher

// loggingServiceMiddleware takes a logger as a dependency
// and returns a service Middleware.
func loggingServiceMiddleware(logger log.Logger) middleware {
	return func(next domain.Publisher) domain.Publisher {
		return loggingMiddleware{
			logger: logger,
			next:   next,
		}
	}
}

type loggingMiddleware struct {
	logger log.Logger
	next   domain.Publisher
}

func (mw loggingMiddleware) Publish(ctx context.Context, event *domain.Event) (err error) {
	defer func() {
		_ = mw.logger.Log("method", "Publish",
			//domain.LogFieldTraceID, traceID,
			//domain.LogFieldSpanID, spanID,
			"event", event,
			"err", err,
		)
	}()
	return mw.next.Publish(ctx, event)
}
