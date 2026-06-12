package memory

import (
	"context"
	"testing"

	"github.com/michaelcardoza/dummy-smtp/internal/core/mail"
)

func TestSave(t *testing.T) {
	repo := New()
	ctx := context.Background()
	msgID := "test-msg-1"

	if err := repo.Save(ctx, &mail.Message{ID: msgID}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got, err := repo.Get(ctx, msgID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got == nil {
		t.Fatalf("got nil, want %v", "Message{}")
	}
	if got.ID != msgID {
		t.Errorf("got %v, want %v", got.ID, msgID)
	}
}
