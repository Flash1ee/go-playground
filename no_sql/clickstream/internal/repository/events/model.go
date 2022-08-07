package events

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"go-clickstream/internal/usecase/events"
)

type event struct {
	ObjectID  primitive.ObjectID `bson:"_id,omitempty"`
	ID        int64              `bson:"id"`
	Body      interface{}        `bson:"event"`
	CreatedAt time.Time          `bson:"created_at"`
	User      interface{}        `bson:"user,omitempty"`
}

func imported(data *events.Event) *event {
	return &event{
		Body: data.Body,
	}
}

func exported(dbEvents []event) []events.Event {
	res := make([]events.Event, len(dbEvents))
	for i, dbEv := range dbEvents {
		switch v := dbEv.Body.(type) {
		case primitive.D:
			res[i].Body = v.Map()
		case primitive.M:
			res[i].Body = v
			res[i].EventID = dbEv.ID
			res[i].Body["createdAt"] = dbEv.CreatedAt
			res[i].Body["user"] = dbEv.User
		default:
			return nil
		}
		res[i].EventID = dbEv.ID
		res[i].Body["createdAt"] = dbEv.CreatedAt
		res[i].Body["user"] = dbEv.User
	}

	return res
}
