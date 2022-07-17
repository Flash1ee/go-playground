package get_event

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"

	"go-clickstream/internal/repository/events"
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
	eventID := ctx.Param("eventID")
	fmt.Println(ctx.Request())
	if eventID == "" {
		return ctx.JSON(http.StatusBadRequest, "empty eventID")
	}
	eventIDAsInt, err := strconv.Atoi(eventID)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest,
			fmt.Sprintf("eventID must be integer, receive eventID = %v", eventID))
	}

	eventFromDB, err := h.usecase.GetEvent(context.Background(), int64(eventIDAsInt))
	if err != nil {
		if err == events.NotFound {
			return ctx.JSON(http.StatusNotFound, struct{}{})
		}
		return err
	}

	return ctx.JSON(http.StatusCreated, eventFromDB)
}
