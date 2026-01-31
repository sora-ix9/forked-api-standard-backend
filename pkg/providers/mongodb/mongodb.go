package mongodb

import (
	"context"
	"fdlp-standard-api/pkg/config"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// MongoDBClient wraps the mongo client and database
type MongoDBClient struct {
	Client *mongo.Client
	DB     *mongo.Database
}

// NewClient initializes and returns a new MongoDB client
func NewClient(cfg *config.Config) *MongoDBClient {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	clientOptions := options.Client().ApplyURI(cfg.MongoDBURI)

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	// Ping the primary
	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		log.Fatalf("Failed to ping MongoDB: %v", err)
	}

	db := client.Database(cfg.MongoDBName)

	return &MongoDBClient{
		Client: client,
		DB:     db,
	}
}

// Disconnect closes the connection to MongoDB
func (m *MongoDBClient) Disconnect(ctx context.Context) error {
	return m.Client.Disconnect(ctx)
}
