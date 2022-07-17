package events

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"go-clickstream/internal/usecase/events"
)

type event struct {
	id        string      `bson:"_id,omitempty"`
	ID        int64       `bson:"id"`
	Body      interface{} `bson:"event"`
	CreatedAt time.Time   `bson:"created_at"`
}

func imported(data *events.Event) *event {
	return &event{
		Body: data.Body,
	}
}

func exported(data *event) *events.Event {
	res := new(events.Event)
	body, ok := data.Body.(primitive.D)
	if !ok {
		return nil
	}

	res.Body = body.Map()
	res.Body["eventID"] = data.ID
	res.Body["createdAt"] = data.CreatedAt

	return res
}
