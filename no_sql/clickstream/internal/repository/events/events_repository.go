package events

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"go-clickstream/internal/usecase/events"
)

var (
	NotFound = errors.New("not found")
)

type Repository struct {
	collection *mongo.Collection
}

func NewRepository(db *mongo.Database, collectionName string) *Repository {
	collection := db.Collection(collectionName)

	return &Repository{
		collection: collection,
	}
}
func (r *Repository) Count(ctx context.Context) (int64, error) {
	res, err := r.collection.CountDocuments(ctx, bson.D{})
	if err != nil {
		return 0, err
	}
	return res, nil
}

func (r *Repository) CreateEvent(ctx context.Context, data *events.Event) (int64, error) {
	cnt, err := r.Count(ctx)
	if err != nil {
		return 0, err
	}
	ev := imported(data)
	ev.ID = cnt + 1

	ev.CreatedAt = time.Now()

	_, err = r.collection.InsertOne(ctx, ev)
	if err != nil {
		return 0, err
	}

	return ev.ID, nil
}

func (r *Repository) GetEvent(ctx context.Context, eventID int64) (interface{}, error) {
	filter := bson.M{"id": eventID}
	res := new(event)

	err := r.collection.FindOne(ctx, filter).Decode(res)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, NotFound
		}
		return 0, err
	}
	return exported(res), nil
}
