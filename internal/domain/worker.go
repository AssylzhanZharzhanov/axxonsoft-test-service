package domain

import "context"

// Worker - represents worker business logic.
type Worker interface {
	Start(ctx context.Context)
}
