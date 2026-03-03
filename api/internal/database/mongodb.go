package database

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDB struct {
	Client   *mongo.Client
	Database *mongo.Database
}

func NewMongoDB(uri string) (*MongoDB, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	clientOptions := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, err
	}

	if err := client.Ping(ctx, nil); err != nil {
		return nil, err
	}

	database := client.Database("flagdeck")

	return &MongoDB{
		Client:   client,
		Database: database,
	}, nil
}

func (m *MongoDB) Close() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return m.Client.Disconnect(ctx)
}

func (m *MongoDB) Ping() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return m.Client.Ping(ctx, nil)
}

func (m *MongoDB) FlagsCollection() *mongo.Collection {
	return m.Database.Collection("flags")
}

func (m *MongoDB) EnvironmentsCollection() *mongo.Collection {
	return m.Database.Collection("environments")
}

func (m *MongoDB) SegmentsCollection() *mongo.Collection {
	return m.Database.Collection("segments")
}

func (m *MongoDB) ExperimentsCollection() *mongo.Collection {
	return m.Database.Collection("experiments")
}

func (m *MongoDB) AuditLogCollection() *mongo.Collection {
	return m.Database.Collection("audit_log")
}

func (m *MongoDB) UsersCollection() *mongo.Collection {
	return m.Database.Collection("users")
}

func (m *MongoDB) APIKeysCollection() *mongo.Collection {
	return m.Database.Collection("api_keys")
}
