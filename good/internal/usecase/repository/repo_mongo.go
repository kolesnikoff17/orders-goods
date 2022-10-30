package repository

import (
	"context"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"good/internal/entity"
	mongodb "good/pkg/mongo"
)

// Good is an implementation of usecase.GoodRepo interface
type Good struct {
	conn *mongodb.MongoDB
}

// New is a contractor for Good
func New(c *mongodb.MongoDB) *Good {
	return &Good{
		conn: c,
	}
}

// GetByID return good with given id
func (m *Good) GetByID(ctx context.Context, id string) (entity.Good, error) {
	collection := m.conn.Client.Database("good_db").Collection("goods")
	var g map[string]interface{}
	err := collection.FindOne(ctx, bson.D{{"_id", id}}).Decode(&g)
	switch {
	case errors.Is(err, mongo.ErrNoDocuments):
		return entity.Good{}, entity.ErrNoID
	case err != nil:
		return entity.Good{}, fmt.Errorf("repo - getbyid: %w", err)
	}
	return entity.Good{ID: id, Data: g}, nil
}

// Create insert new good in collection and return its id
func (m *Good) Create(ctx context.Context, good entity.Good) (string, error) {
	collection := m.conn.Client.Database("good_db").Collection("goods")
	doc, _ := bson.Marshal(good.Data)
	res, err := collection.InsertOne(ctx, doc)
	if err != nil {
		return "", fmt.Errorf("repo - create: %w", err)
	}
	return res.InsertedID.(string), nil
}

// Update updates good with given id
func (m *Good) Update(ctx context.Context, good entity.Good) error {
	collection := m.conn.Client.Database("good_db").Collection("goods")
	doc, _ := bson.Marshal(good.Data)
	_, err := collection.UpdateByID(ctx, good.ID, doc)
	if err != nil {
		return fmt.Errorf("repo - update: %w", err)
	}
	return nil
}

// Delete remove good with given id
func (m *Good) Delete(ctx context.Context, id string) error {
	collection := m.conn.Client.Database("good_db").Collection("goods")
	_, err := collection.DeleteOne(ctx, bson.D{{"_id", id}})
	if err != nil {
		return fmt.Errorf("repo - delete: %w", err)
	}
	return nil
}
