package send_event

import (
	"context"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"

	"go-clickstream/internal/model/dto"
	us "go-clickstream/internal/usecase/events"
)

type Handler struct {
	usecase usecase
}

func NewHandler(usecase usecase) *Handler {
	return &Handler{
		usecase: usecase,
	}
}

func (h *Handler) Handle(ctx echo.Context) error {
	request := new(dto.RequestSendEvent)

	if err := ctx.Bind(request); err != nil {
		return err
	}
	fillRequestParams(request, ctx.Request())

	id, err := h.usecase.CreateEvent(context.Background(), &us.Event{Body: *request})
	if err != nil {
		return err
	}

	response := struct {
		ID int64 `json:"eventID"`
	}{ID: id}

	return ctx.JSON(http.StatusCreated, response)
}

func fillRequestParams(data *dto.RequestSendEvent, req *http.Request) {
	(*data)["method"] = req.Method
	(*data)["headers"] = headerToArray(req.Header)
	(*data)["content-type"] = req.Cookies()
}

func headerToArray(header http.Header) (res []string) {
	for name, values := range header {
		for _, value := range values {
			res = append(res, fmt.Sprintf("%s: %s", name, value))
		}
	}
	return
}
