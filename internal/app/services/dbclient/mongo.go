package dbclient

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"task/internal/app/services/configmanager"
)

type MongoDBClient struct {
	config *configmanager.Config
	client *mongo.Client
}

func NewMongoDBClient(config *configmanager.Config) *MongoDBClient {
	return &MongoDBClient{
		config: config,
	}
}

func (c *MongoDBClient) Connect(ctx context.Context) (*mongo.Client, error) {
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017/")
	var err error
	c.client, err = mongo.Connect(ctx, clientOptions)
	return c.client, err
}

func (c *MongoDBClient) Disconnect(ctx context.Context) {
	if err := c.client.Disconnect(ctx); err != nil {
		log.Fatal(err)
	}
}
