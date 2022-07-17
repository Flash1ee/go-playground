package events

type Event struct {
	Body map[string]interface{} `json:"event"`
}
