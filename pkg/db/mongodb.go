package db

import (
	"context"
	"fdlp-standard-api/pkg/config"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type MongoDB struct {
	Client  *mongo.Client
	DB      *mongo.Database
	Context context.Context
	timeout time.Duration
}

func NewMongoDB(timeoutSeconds int) *MongoDB {
	return NewMongoDBWithConfig(config.New(), timeoutSeconds)
}

func NewMongoDBWithConfig(cfg *config.Config, timeoutSeconds int) *MongoDB {
	timeout := time.Duration(timeoutSeconds) * time.Second
	if timeout <= 0 {
		timeout = 10 * time.Second
	}

	baseCtx := context.Background()
	ctx, cancel := context.WithTimeout(baseCtx, timeout)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(cfg.MongoDBURI))
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		log.Fatalf("Failed to ping MongoDB: %v", err)
	}

	return &MongoDB{
		Client:  client,
		DB:      client.Database(cfg.MongoDBDb),
		Context: baseCtx,
		timeout: timeout,
	}
}

func (m *MongoDB) Disconnect() {
	if m == nil || m.Client == nil {
		return
	}

	ctx := m.Context
	if ctx == nil {
		ctx = context.Background()
	}

	timeout := m.timeout
	if timeout <= 0 {
		timeout = 10 * time.Second
	}

	disconnectCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	if err := m.Client.Disconnect(disconnectCtx); err != nil {
		log.Printf("Failed to disconnect MongoDB: %v", err)
	}
}
