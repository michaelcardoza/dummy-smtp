package mongo

import (
	"context"

	"github.com/michaelcardoza/dummy-smtp/internal/core/mail"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type Storage struct {
	client     *mongo.Client
	db         *mongo.Database
	collection *mongo.Collection
}

func NewStorage(client *mongo.Client) *Storage {
	db := client.Database("dummysmtp")
	return &Storage{
		client:     client,
		db:         db,
		collection: db.Collection("messages"),
	}
}

func (s *Storage) Save(ctx context.Context, message *mail.Message) error {
	if _, err := s.collection.InsertOne(ctx, message); err != nil {
		return err
	}
	return nil
}

func (s *Storage) List(ctx context.Context) ([]*mail.Message, error) {
	opts := options.Find().SetSort(bson.D{{Key: "createdAt", Value: -1}})
	cursor, err := s.collection.Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	messages := make([]*mail.Message, 0)
	if err = cursor.All(ctx, &messages); err != nil {
		return nil, err
	}

	return messages, nil
}

func (s *Storage) Get(ctx context.Context, id string) (*mail.Message, error) {
	var message mail.Message

	filter := bson.M{"id": id}
	err := s.collection.FindOne(ctx, filter).Decode(&message)
	if err != nil {
		return nil, err
	}

	return &message, nil
}

func (s *Storage) DeleteByID(ctx context.Context, id string) error {
	filter := bson.M{"id": id}
	_, err := s.collection.DeleteOne(ctx, filter)
	return err
}

func (s *Storage) DeleteAll(ctx context.Context) error {
	_, err := s.collection.DeleteMany(ctx, bson.D{})
	return err
}

func (s *Storage) Close(ctx context.Context) error {
	return s.client.Disconnect(ctx)
}
