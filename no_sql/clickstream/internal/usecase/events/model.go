package events

type Event struct {
	Body    map[string]interface{} `json:"event"`
	EventID int64                  `json:"eventID"`
}
