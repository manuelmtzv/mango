package db

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DB struct {
	Client *mongo.Client
	DB     *mongo.Database
}

func New(uri, dbName string) (*DB, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	clientOptions := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to mongo: %w", err)
	}

	if err := client.Ping(ctx, nil); err != nil {
		return nil, fmt.Errorf("failed to ping mongo: %w", err)
	}

	log.Println("Connected to MongoDB")
	return &DB{
		Client: client,
		DB:     client.Database(dbName),
	}, nil
}

func (d *DB) Close() {
	if err := d.Client.Disconnect(context.Background()); err != nil {
		log.Printf("Error closing mongo connection: %v", err)
	}
}
