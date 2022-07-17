package events

import (
	"context"
)

type repository interface {
	CreateEvent(ctx context.Context, data *Event) (int64, error)
	GetEvent(ctx context.Context, eventID int64) (interface{}, error)
}
