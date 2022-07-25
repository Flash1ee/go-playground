package dto

const (
	NotFoundUserID = -1
	InvalidUserID  = -2
)

type RequestSendEvent map[string]interface{}

func (m *RequestSendEvent) GetUserID() int64 {
	userID := (*m)["userID"]
	if userID == "" {
		return NotFoundUserID
	}
	userIDfloat, ok := userID.(float64)
	if !ok || userIDfloat < 1 {
		return InvalidUserID
	}

	if !isInt(userIDfloat) {
		return InvalidUserID
	}

	return int64(userIDfloat)
}

func isInt(val float64) bool {
	return val == float64(int64(val))
}
