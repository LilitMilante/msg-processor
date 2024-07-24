package service

import (
	"context"
	"time"

	"msg-processor/internal/entity"

	"github.com/google/uuid"
)

type Repository interface {
	CreateMsg(ctx context.Context, msg entity.Msg) error
}

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) CreateMsg(ctx context.Context, text string) (entity.Msg, error) {
	msg := entity.Msg{
		ID:        uuid.New(),
		Text:      text,
		State:     entity.StateUnprocessed,
		CreatedAt: time.Now(),
	}

	err := s.repo.CreateMsg(ctx, msg)
	if err != nil {
		return entity.Msg{}, err
	}

	return msg, nil
}
