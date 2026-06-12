package mail

import (
	"context"
	"errors"
	"testing"
)

type fakeRepo struct {
	saved   map[string]*Message
	saveErr error
}

func newFakeRepo() *fakeRepo {
	return &fakeRepo{
		saved: make(map[string]*Message),
	}
}

func (r *fakeRepo) Save(_ context.Context, message *Message) error {
	if r.saveErr != nil {
		return r.saveErr
	}
	r.saved[message.ID] = message
	return nil
}

func (r *fakeRepo) List(_ context.Context) ([]*Message, error) {
	messages := make([]*Message, 0, len(r.saved))
	for _, message := range r.saved {
		messages = append(messages, message)
	}
	return messages, nil
}

func (r *fakeRepo) Get(_ context.Context, id string) (*Message, error) {
	message, ok := r.saved[id]
	if !ok {
		return nil, ErrNotFound
	}
	return message, nil
}

func (r *fakeRepo) DeleteById(_ context.Context, id string) error {
	if _, ok := r.saved[id]; !ok {
		return ErrNotFound
	}
	delete(r.saved, id)
	return nil
}

func TestService_Capture(t *testing.T) {
	repo := newFakeRepo()
	svc := NewService(repo)

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
	if _, ok := repo.saved[message.ID]; !ok {
		t.Error("message not saved")
	}
}

func TestService_CaptureNoRecipients(t *testing.T) {
	repo := newFakeRepo()
	svc := NewService(repo)
	_, err := svc.Capture(context.Background(), CaptureParams{From: "alice@example.com"})
	if !errors.Is(err, ErrNoRecipients) {
		t.Errorf("expected ErrNoRecipients, got: %v", err)
	}
}

func TestService_CaptureSaveError(t *testing.T) {
	repo := newFakeRepo()
	repo.saveErr = errors.New("disk full")
	svc := NewService(repo)
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
