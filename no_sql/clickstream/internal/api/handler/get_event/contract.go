package get_event

import (
	"context"

	"go-clickstream/internal/usecase/events"
)

type usecase interface {
	GetEvent(ctx context.Context, eventID int64) (*events.Event, error)
}
