package memory

import (
	"context"
	"sync"

	"github.com/michaelcardoza/dummy-smtp/internal/core/mail"
)

type Storage struct {
	mu       sync.RWMutex
	messages []*mail.Message
}

func NewStorage() *Storage {
	return &Storage{}
}

func (s *Storage) Save(_ context.Context, message *mail.Message) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.messages = append(s.messages, message)
	return nil
}

func (s *Storage) Get(_ context.Context, id string) (*mail.Message, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, message := range s.messages {
		if message.ID == id {
			return message, nil
		}
	}
	return nil, mail.ErrNotFound
}

func (s *Storage) List(_ context.Context) ([]*mail.Message, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	messages := make([]*mail.Message, len(s.messages))
	for i, message := range s.messages {
		messages[len(s.messages)-1-i] = message
	}
	return messages, nil
}

func (s *Storage) DeleteByID(_ context.Context, id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for i, message := range s.messages {
		if message.ID == id {
			s.messages = append(s.messages[:i], s.messages[i+1:]...)
			return nil
		}
	}
	return mail.ErrNotFound
}

func (s *Storage) DeleteAll(_ context.Context) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.messages = nil
	return nil
}

func (s *Storage) Close(_ context.Context) error {
	return nil
}
