package events

import (
	"context"
)

type repository interface {
	CreateEvent(ctx context.Context, data *Event) (int64, error)
	GetEvent(ctx context.Context, eventID int64) ([]Event, error)
	UserExist(ctx context.Context, userID int64) (bool, error)
	UploadEvents(ctx context.Context, userID int64) ([]Event, error)
}
