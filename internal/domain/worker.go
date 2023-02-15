package domain

import "context"

type Worker interface {
	Start(ctx context.Context)
}
