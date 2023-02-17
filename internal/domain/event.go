package domain

import (
	"context"
	"fmt"
)

// Event - represents event
type Event struct {
	TaskID    TaskID
	CreatedAt int64
}

// Validate - validates struct.
func (e Event) Validate() error {
	if e.TaskID <= 0 {
		return fmt.Errorf("task id required")
	}

	return nil
}

// EventService -  provides access to a business logic.
type EventService interface {
	// RegisterEvent - creates a new event and publish it to queue
	//
	RegisterEvent(ctx context.Context, event *Event) error
}
