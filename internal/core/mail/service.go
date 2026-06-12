package mail

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type CaptureParams struct {
	From        string
	To          []string
	Subject     string
	TextBody    string
	HTMLBody    string
	Raw         string
	Attachments []Attachment
}

type Service struct {
	repo MessageRepository
	pub  Publisher
}

func NewService(repo MessageRepository, pub Publisher) *Service {
	return &Service{repo: repo, pub: pub}
}

func (s *Service) Capture(ctx context.Context, params CaptureParams) (*Message, error) {
	message := &Message{
		ID:          uuid.New().String(),
		From:        params.From,
		To:          params.To,
		Subject:     params.Subject,
		TextBody:    params.TextBody,
		HTMLBody:    params.HTMLBody,
		Raw:         params.Raw,
		Attachments: params.Attachments,
		CreatedAt:   time.Now(),
	}
	if !message.HasRecipients() {
		return nil, ErrNoRecipients
	}

	if err := s.repo.Save(ctx, message); err != nil {
		return nil, fmt.Errorf("save message: %w", err)
	}

	s.pub.Publish(message)

	return message, nil
}

func (s *Service) List(ctx context.Context) ([]*Message, error) {
	return s.repo.List(ctx)
}

func (s *Service) Get(ctx context.Context, id string) (*Message, error) {
	return s.repo.Get(ctx, id)
}

func (s *Service) DeleteByID(ctx context.Context, id string) error {
	return s.repo.DeleteByID(ctx, id)
}

func (s *Service) DeleteAll(ctx context.Context) error {
	return s.repo.DeleteAll(ctx)
}
