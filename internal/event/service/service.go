package service

import (
	"context"

	"github.com/AssylzhanZharzhanov/axxonsoft-test-service/internal/domain"
)

type service struct {
}

func NewService() domain.EventService {
	return &service{}
}

func (s service) RegisterEvent(ctx context.Context) {
	//TODO implement me
	panic("implement me")
}
