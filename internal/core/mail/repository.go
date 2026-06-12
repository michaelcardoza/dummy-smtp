package mail

import "context"

type MessageRepository interface {
	Save(ctx context.Context, message *Message) error
	Get(ctx context.Context, id string) (*Message, error)
	List(ctx context.Context) ([]*Message, error)
	DeleteByID(ctx context.Context, id string) error
	DeleteAll(ctx context.Context) error
}
