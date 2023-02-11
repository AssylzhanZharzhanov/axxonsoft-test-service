package domain

import "context"

type Event struct {
}

type EventService interface {
	RegisterEvent(ctx context.Context)
}
