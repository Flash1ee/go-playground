package send_event

import (
	"context"

	"go-clickstream/internal/usecase/events"
)

type usecase interface {
	CreateEvent(ctx context.Context, data *events.Event) (int64, error)
}
