package domain

import "context"

// Consumer - represents consumer
type Consumer interface {
	Consume(ctx context.Context) error
}
