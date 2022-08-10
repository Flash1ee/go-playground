package events

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/gommon/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"go-clickstream/internal/usecase/events"
)

var (
	NotFound       = errors.New("not found")
	NotFoundEvents = errors.New("events not found")
)

type Repository struct {
	eventsCollection *mongo.Collection
	userCollection   *mongo.Collection
	client           *mongo.Client
}

func NewRepository(c *mongo.Client, db *mongo.Database, eventCollection string, userCollection string) *Repository {
	evCollection := db.Collection(eventCollection)
	usCollection := db.Collection(userCollection)

	return &Repository{
		eventsCollection: evCollection,
		userCollection:   usCollection,
		client:           c,
	}
}
func (r *Repository) count(ctx context.Context) (int64, error) {
	filter := primitive.D{
		{"id", primitive.D{{"$exists", "true"}}},
	}
	opts := options.Find().SetSort(bson.D{{"id", -1}}).SetLimit(1)
	var data []struct {
		ID int64 `json:"id"`
	}
	cur, err := r.eventsCollection.Find(ctx, filter, opts)
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
	cnt, err := r.count(ctx)
	if err != nil && err != NotFound {
		return 0, err
	}
	ev := imported(data)
	ev.ID = cnt + 1

	ev.CreatedAt = time.Now()

	_, err = r.eventsCollection.InsertOne(ctx, ev)
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

func (r *Repository) UserExist(ctx context.Context, userID int64) (bool, error) {
	filter := primitive.D{
		{Key: "userID", Value: primitive.D{{Key: "$eq", Value: userID}}},
	}
	err := r.userCollection.FindOne(ctx, filter).Err()
	if err == nil {
		return true, nil
	}
	if err == mongo.ErrNoDocuments {
		return false, NotFound
	}
	return false, err
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
	cursor, err := r.eventsCollection.Aggregate(ctx, mongo.Pipeline{lookupStage, matchStage})
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

// db.events.find({"event.userID":  2})
// 1. Выгрузить данные по userID
// 2. Обновить userID на UUID
// 3. Вернуть выгруженные данные
/*
db.events.updateMany({"event.userID": 5}, [{
    $set:
        {
            "event.userID": {
                $function: {
                    body: "function () {return UUID().toString().split('\"')[1];}",
                    args: [],
                    lang: "js",
                }
            },
        }
}])
*/
func (r *Repository) UploadEvents(ctx context.Context, userID int64) ([]events.Event, error) {
	findFilter := bson.D{{Key: "event.userID", Value: userID}}
	//updateFilter := bson.D{
	//	{Key: "$set", Value: bson.D{
	//		{Key: "event.userID", Value: bson.D{
	//			{Key: "$function", Value: bson.D{
	//				{Key: "body", Value: "function () {return UUID().toString().split('\"')[1];}"},
	//				{Key: "args", Value: "[]"},
	//				{Key: "lang", Value: "js"},
	//			}},
	//		}}}}}
	session, err := r.client.StartSession()
	if err != nil {
		return nil, err
	}
	defer session.EndSession(ctx)

	res := make([]event, 0, 2)

	err = mongo.WithSession(ctx, session, func(sc mongo.SessionContext) error {
		if err = session.StartTransaction(); err != nil {
			return err
		}

		cur, err := r.eventsCollection.Find(sc, findFilter)
		if err != nil {
			_ = session.AbortTransaction(sc)
			return err
		}

		for cur.Next(sc) {
			ev := event{}
			err = cur.Decode(&ev)
			if err != nil {
				_ = session.AbortTransaction(sc)
				return err
			}
			newID := uuid.New()
			if err != nil {
				_ = session.AbortTransaction(sc)
				return err
			}
			body, ok := ev.Body.(primitive.D)
			if !ok {
				_ = session.AbortTransaction(sc)
				return fmt.Errorf("can not convert event body to map[string]interface{}")
			}
			bodyMap := body.Map()
			bodyMap["userID"] = newID
			ev.Body = bodyMap

			res = append(res, ev)
		}
		if len(res) == 0 {
			_ = session.AbortTransaction(sc)
			return NotFoundEvents
		}

		for _, ev := range res {
			body, ok := ev.Body.(primitive.M)
			if !ok {
				_ = session.AbortTransaction(sc)
				return fmt.Errorf("can not convert event body to map[string]interface{} in upd loop")
			}
			log.Info(ev.ObjectID.String())
			result, err := r.eventsCollection.UpdateByID(
				sc, ev.ObjectID, bson.D{
					{Key: "$set", Value: bson.D{
						{"event.userID", body["userID"]}},
					},
				})
			if err != nil {
				_ = session.AbortTransaction(sc)
				return err
			}
			cnt := result.ModifiedCount
			if cnt != 1 {
				_ = session.AbortTransaction(sc)
				return fmt.Errorf("ModifiedCount not eq 1 - error")
			}
		}
		if err = session.CommitTransaction(sc); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return exported(res), nil
}
