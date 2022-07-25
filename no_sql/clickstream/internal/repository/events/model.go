package events

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"go-clickstream/internal/usecase/events"
)

type event struct {
	id        string      `bson:"_id,omitempty"`
	ID        int         `bson:"id"`
	Body      interface{} `bson:"event"`
	CreatedAt time.Time   `bson:"created_at"`
	User      interface{} `bson:"user"`
}

func imported(data *events.Event) *event {
	return &event{
		Body: data.Body,
	}
}

func exported(dbEvents []event) []events.Event {
	res := make([]events.Event, len(dbEvents))
	for i, dbEv := range dbEvents {
		body, ok := dbEv.Body.(primitive.D)
		if !ok {
			return nil
		}

		res[i].Body = body.Map()
		res[i].EventID = int64(dbEv.ID)
		res[i].Body["createdAt"] = dbEv.CreatedAt
		res[i].Body["user"] = dbEv.User

	}

	return res
}
