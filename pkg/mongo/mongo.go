package mongodb

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// MongoDB keeps connection to db
type MongoDB struct {
	Client *mongo.Client
}

// New is a contractor for MongoDB
func New(uri string) (*MongoDB, error) {
	ctx := context.Background()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		return nil, err
	}
	return &MongoDB{Client: client}, nil
}

// Close disconnects db
func (conn *MongoDB) Close() {
	ctx := context.Background()
	_ = conn.Client.Disconnect(ctx)
}
