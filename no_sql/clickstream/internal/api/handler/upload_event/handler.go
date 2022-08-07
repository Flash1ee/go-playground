package upload_event

import (
	"context"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"

	"go-clickstream/internal/model/dto"
	events2 "go-clickstream/internal/repository/events"
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

	events, err := h.usecase.UploadEvents(context.Background(), userID)
	if err != nil {
		if err == events2.NotFound {
			return ctx.JSON(http.StatusNotFound, dto.ErrorResponse{
				Message: "userID not found in request body",
			})
		}
		if err == events2.NotFoundEvents {
			return ctx.JSON(http.StatusNotFound, dto.ErrorResponse{
				Message: "events by this userID not found",
			})
		}
		log.Error(fmt.Errorf("GetEventsByUserID error: %w", err))
		return err
	}

	return ctx.JSON(http.StatusOK, present(events))
}
