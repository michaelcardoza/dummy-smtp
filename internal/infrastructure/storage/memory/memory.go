package memory

import (
	"context"
	"sync"

	"github.com/michaelcardoza/dummy-smtp/internal/core/mail"
)

type Repository struct {
	mu       sync.RWMutex
	messages []*mail.Message
}

func New() *Repository {
	return &Repository{}
}

func (r *Repository) Save(_ context.Context, message *mail.Message) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.messages = append(r.messages, message)
	return nil
}

func (r *Repository) Get(_ context.Context, id string) (*mail.Message, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, message := range r.messages {
		if message.ID == id {
			return message, nil
		}
	}
	return nil, mail.ErrNotFound
}

func (r *Repository) List(_ context.Context) ([]*mail.Message, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	messages := make([]*mail.Message, len(r.messages))
	for i, message := range r.messages {
		messages[len(r.messages)-1-i] = message
	}
	return messages, nil
}

func (r *Repository) DeleteByID(_ context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	for i, message := range r.messages {
		if message.ID == id {
			r.messages = append(r.messages[:i], r.messages[i+1:]...)
			return nil
		}
	}
	return mail.ErrNotFound
}

func (r *Repository) DeleteAll(_ context.Context) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.messages = []*mail.Message{}
	return nil
}
