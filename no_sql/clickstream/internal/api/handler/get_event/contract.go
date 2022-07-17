package get_event

import (
	"context"
)

type usecase interface {
	GetEvent(ctx context.Context, eventID int64) (interface{}, error)
}
