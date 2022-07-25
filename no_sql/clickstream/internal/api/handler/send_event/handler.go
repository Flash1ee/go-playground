package send_event

import (
	"context"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"

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
	userID := request.GetUserID()
	if userID == dto.InvalidUserID {
		return ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Message: "userID is invalid, valid only int64 > 0",
		})
	}

	if userID == dto.NotFoundUserID {
		return ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Message: "userID not found in request body",
		})

	}

	fmt.Println(userID)
	fillRequestParams(request, ctx.Request())

	id, err := h.usecase.CreateEvent(context.Background(), &us.Event{Body: *request})
	if err != nil {
		log.Error(fmt.Errorf("CreateEvent error: %w", err))
		// вернет   "message": "Internal Server Error"
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
