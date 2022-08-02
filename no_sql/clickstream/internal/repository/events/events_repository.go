package events

import (
	"context"
	"errors"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

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
	filter := primitive.D{
		{"id", primitive.D{{"$exists", "true"}}},
	}
	opts := options.Find().SetSort(bson.D{{"id", -1}}).SetLimit(1)
	var data []struct {
		ID int64 `json:"id"`
	}
	cur, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return 0, err
	}
	err = cur.All(ctx, &data)
	if err != nil {
		return 0, err
	}
	if len(data) == 0 {
		return 0, NotFound
	}

	return data[0].ID, nil
}

func (r *Repository) CreateEvent(ctx context.Context, data *events.Event) (int64, error) {
	cnt, err := r.Count(ctx)
	if err != nil && err != NotFound {
		return 0, err
	}
	ev := imported(data)
	ev.ID = cnt + 1

	ev.CreatedAt = time.Now()

	_, err = r.collection.InsertOne(ctx, ev)
	if err != nil {
		we, ok := err.(mongo.WriteException)
		if ok {
			for _, r := range we.WriteErrors {
				err = fmt.Errorf("%w", r)
			}

		}
		return 0, err
	}
	return ev.ID, nil
}

// db.events.aggregate([
//    {
//        "$match":
//            {
//                "id": {
//                    "$eq": 1
//                },
//            },
//    },
//    {
//        $lookup: {
//            from: "users",
//            localField: "event.userID",
//            foreignField: "id",
//            as: "user"
//        }
//    },
//])
func (r *Repository) GetEvent(ctx context.Context, eventID int64) ([]events.Event, error) {
	matchStage := bson.D{{"$match", bson.D{{"id", eventID}}}}
	lookupStage := bson.D{{"$lookup", bson.D{
		{"from", "users"},
		{"localField", "event.userID"},
		{"foreignField", "userID"},
		{"as", "user"}}}}
	var res []event
	cursor, err := r.collection.Aggregate(ctx, mongo.Pipeline{lookupStage, matchStage})
	if err != nil {
		return nil, err
	}
	err = cursor.All(ctx, &res)
	if err != nil {
		return nil, err
	}
	if len(res) == 0 {
		return nil, NotFound
	}
	return exported(res), nil
}
