package upload_event

import (
	"time"

	events2 "go-clickstream/internal/usecase/events"
)

func present(events []events2.Event) []events2.Event {
	response := make([]events2.Event, 0, len(events))

	for _, val := range events {
		evTime, ok := val.Body["created_at"]
		if !ok {
			evTime = time.UnixDate
		}
		response = append(response, events2.Event{
			EventID: val.EventID,
			Body: map[string]interface{}{
				"created_at": evTime,
			},
		})
	}
	return response
}
