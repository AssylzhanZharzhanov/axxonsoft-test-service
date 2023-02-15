package domain

import (
	"context"
	"fmt"
)

type Event struct {
	TaskID    TaskID
	CreatedAt int64
}

func (e Event) Validate() error {
	if e.TaskID <= 0 {
		return fmt.Errorf("task id required")
	}

	return nil
}

type EventService interface {
	RegisterEvent(ctx context.Context, event *Event) error
}
