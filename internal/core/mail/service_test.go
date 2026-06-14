package mail

import (
	"context"
	"errors"
	"testing"
)

type fakeStorage struct {
	saved   map[string]*Message
	saveErr error
}

func newFakeStorage() *fakeStorage {
	return &fakeStorage{
		saved: make(map[string]*Message),
	}
}

func (s *fakeStorage) Save(_ context.Context, message *Message) error {
	if s.saveErr != nil {
		return s.saveErr
	}
	s.saved[message.ID] = message
	return nil
}

func (s *fakeStorage) List(_ context.Context) ([]*Message, error) {
	messages := make([]*Message, 0, len(s.saved))
	for _, message := range s.saved {
		messages = append(messages, message)
	}
	return messages, nil
}

func (s *fakeStorage) Get(_ context.Context, id string) (*Message, error) {
	message, ok := s.saved[id]
	if !ok {
		return nil, ErrNotFound
	}
	return message, nil
}

func (s *fakeStorage) DeleteByID(_ context.Context, id string) error {
	if _, ok := s.saved[id]; !ok {
		return ErrNotFound
	}
	delete(s.saved, id)
	return nil
}

func (s *fakeStorage) DeleteAll(_ context.Context) error {
	s.saved = nil
	return nil
}

func (s *fakeStorage) Close() error {
	return nil
}

func TestService_Capture(t *testing.T) {
	storage := newFakeStorage()
	svc := NewService(storage, nil)

	message, err := svc.Capture(context.Background(), CaptureParams{
		From:     "alice@example.com",
		To:       []string{"bob@example.com"},
		Subject:  "Hello 👋",
		TextBody: "plain",
		HTMLBody: "<p>html</p>",
		Raw:      "row original",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if message.ID == "" {
		t.Error("ID was not generated")
	}
	if message.CreatedAt.IsZero() {
		t.Error("CreatedAt was not set")
	}
	if message.From != "alice@example.com" {
		t.Fatalf("got %q, want %s", message.From, "alice@example.com")
	}
	if message.HTMLBody != "<p>html</p>" {
		t.Fatalf("got %q, want %s", message.HTMLBody, "<p>html</p>")
	}
	if _, ok := storage.saved[message.ID]; !ok {
		t.Error("message not saved")
	}
}

func TestService_CaptureNoRecipients(t *testing.T) {
	storage := newFakeStorage()
	svc := NewService(storage, nil)
	_, err := svc.Capture(context.Background(), CaptureParams{From: "alice@example.com"})
	if !errors.Is(err, ErrNoRecipients) {
		t.Errorf("expected ErrNoRecipients, got: %v", err)
	}
}

func TestService_CaptureSaveError(t *testing.T) {
	storage := newFakeStorage()
	storage.saveErr = errors.New("disk full")
	svc := NewService(storage, nil)
	_, err := svc.Capture(context.Background(), CaptureParams{
		From: "alice@example.com",
		To:   []string{"bob@example.com"},
	})
	if err == nil {
		t.Fatal("expected error when Save fails")
	}
}

func TestMessage_HasRecipients(t *testing.T) {
	if (&Message{}).HasRecipients() {
		t.Error("empty To should report no recipients")
	}
	if !(&Message{To: []string{"bob@example.com"}}).HasRecipients() {
		t.Error("non-empty To should report recipients")
	}
}
