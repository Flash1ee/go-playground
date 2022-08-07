package events

import (
	"context"
)

type EventsUsecase struct {
	repository repository
}

func NewEventsUsecase(repository repository) *EventsUsecase {
	return &EventsUsecase{
		repository: repository,
	}
}

func (u *EventsUsecase) CreateEvent(ctx context.Context, data *Event) (int64, error) {
	res, err := u.repository.CreateEvent(ctx, data)

	return res, err
}

func (u *EventsUsecase) GetEvent(ctx context.Context, eventID int64) (*Event, error) {
	res, err := u.repository.GetEvent(ctx, eventID)
	if err != nil {
		return nil, err
	}

	return &res[0], nil
}

func (u *EventsUsecase) UploadEvents(ctx context.Context, userID int64) ([]Event, error) {
	if ok, err := u.repository.UserExist(ctx, userID); !ok {
		return nil, err
	}
	res, err := u.repository.UploadEvents(ctx, userID)
	if err != nil {
		return nil, err
	}

	return res, nil
}
