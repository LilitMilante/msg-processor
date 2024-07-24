package service

import (
	"context"
	"time"

	"msg-processor/internal/entity"

	"github.com/gofrs/uuid"
)

type Repository interface {
	CreateMsg(ctx context.Context, msg entity.Msg) error
	UnprocessedMsgs(ctx context.Context) ([]entity.Msg, error)
	UpdateMsgsState(ctx context.Context, state entity.State, msgs ...entity.Msg) error
}

type Producer interface {
	SendMsg(ctx context.Context, msgs ...entity.Msg) error
}

type Service struct {
	repo     Repository
	producer Producer
}

func NewService(repo Repository, producer Producer) *Service {
	return &Service{repo: repo, producer: producer}
}

func (s *Service) CreateMsg(ctx context.Context, text string) (entity.Msg, error) {
	msg := entity.Msg{
		ID:        uuid.Must(uuid.NewV4()),
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

func (s *Service) SendMsgToKafka(ctx context.Context) error {
	msgs, err := s.repo.UnprocessedMsgs(ctx)
	if err != nil {
		return err
	}

	err = s.producer.SendMsg(ctx, msgs...)
	if err != nil {
		return err
	}

	err = s.repo.UpdateMsgsState(ctx, entity.StateProcessing, msgs...)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) ProcessMsg(ctx context.Context, msg entity.Msg) error {
	return s.repo.UpdateMsgsState(ctx, entity.StateProcessed, msg)
}
