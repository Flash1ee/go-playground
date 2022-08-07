package upload_event

import (
	"context"

	"go-clickstream/internal/usecase/events"
)

type usecase interface {
	UploadEvents(ctx context.Context, userID int64) ([]events.Event, error)
}
