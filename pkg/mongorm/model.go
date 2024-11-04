package mongorm

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Model struct {
	ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time          `json:"updated_at" bson:"updated_at"`
}

func (m *Model) GetID() string {
	return m.ID.Hex()
}

func (m *Model) Create(ctx context.Context, db *mongo.Database, collectionName string, model interface{}, opts ...*options.InsertOneOptions) error {
	collection := db.Collection(collectionName)

	m.CreatedAt = time.Now()
	m.UpdatedAt = time.Now()

	res, err := collection.InsertOne(ctx, model, opts...)
	if err != nil {
		return err
	}

	m.ID = res.InsertedID.(primitive.ObjectID)
	return nil
}

func (m *Model) Read(ctx context.Context, db *mongo.Database, collectionName string, filter interface{}, result interface{}, opts ...*options.FindOneOptions) error {
	collection := db.Collection(collectionName)

	err := collection.FindOne(ctx, filter, opts...).Decode(result)
	if err != nil {
		if errors.Is(err, bson.ErrNilRegistry) {
			return nil
		}
		return err
	}

	return nil
}

func List(ctx context.Context, db *mongo.Database, collectionName string, filter interface{}, results interface{}, opts ...*options.FindOptions) error {
	collection := db.Collection(collectionName)

	cursor, err := collection.Find(ctx, filter, opts...)
	if err != nil {
		return err
	}

	if err = cursor.All(ctx, results); err != nil {
		return err
	}

	return nil
}

func (m *Model) Update(ctx context.Context, db *mongo.Database, collectionName string, filter interface{}, update interface{}, opts ...*options.FindOneAndUpdateOptions) error {
	collection := db.Collection(collectionName)

	m.UpdatedAt = time.Now()
	update.(bson.M)["$set"].(bson.M)["updated_at"] = m.UpdatedAt

	err := collection.FindOneAndUpdate(ctx, filter, update, opts...).Decode(&m)
	if err != nil && !errors.Is(err, mongo.ErrNoDocuments) {
		return err
	}

	if m.CreatedAt == (time.Time{}) {
		m.CreatedAt = time.Now()
		_, err := collection.UpdateByID(ctx, m.ID, bson.M{"$set": bson.M{"created_at": m.CreatedAt, "updated_at": m.UpdatedAt}})
		if err != nil {
			return err
		}
	}

	return nil
}

func Aggregate(ctx context.Context, db *mongo.Database, collectionName string, pipeline mongo.Pipeline, results interface{}, opts ...*options.AggregateOptions) error {
	collection := db.Collection(collectionName)

	cursor, err := collection.Aggregate(ctx, pipeline, opts...)
	if err != nil {
		return err
	}

	if err = cursor.All(ctx, results); err != nil {
		return err
	}
	return nil
}

func (m *Model) Delete(ctx context.Context, db *mongo.Database, collectionName string, filter interface{}, opts ...*options.DeleteOptions) error {
	collection := db.Collection(collectionName)
	_, err := collection.DeleteOne(ctx, filter, opts...)
	if err != nil {
		return err
	}

	return nil
}
