package service

import (
	"context"
	"testing"
	"time"

	"github.com/AssylzhanZharzhanov/axxonsoft-test-service/internal/domain"
	"github.com/AssylzhanZharzhanov/axxonsoft-test-service/internal/mocks"

	"github.com/golang/mock/gomock"
)

func TestService_RegisterEvent(t *testing.T) {

	var (
		ctx        = context.Background()
		validEvent = &domain.Event{
			TaskID:    1,
			CreatedAt: time.Now().Unix(),
		}
	)

	// Setup mocks.
	stubCtrl := gomock.NewController(t)
	defer stubCtrl.Finish()

	publisherStub := mocks.NewMockPublisher(stubCtrl)
	service := newBasicService(publisherStub)

	publisherStub.EXPECT().
		Publish(ctx, validEvent).
		Return(nil).
		AnyTimes()

	tests := []struct {
		name        string
		event       *domain.Event
		expectError bool
	}{
		{
			name:        "Success: event registered",
			event:       validEvent,
			expectError: false,
		},
		{
			name: "Fail: invalid task id",
			event: &domain.Event{
				TaskID:    0,
				CreatedAt: time.Now().Unix(),
			},
			expectError: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			event := test.event
			err := service.RegisterEvent(ctx, event)
			if !test.expectError {
				if err != nil {
					t.Errorf("unexpected error: %s", err)
				}
			} else {
				if err == nil {
					t.Errorf("unexpected error got nothing")
				}
			}
		})
	}
}
