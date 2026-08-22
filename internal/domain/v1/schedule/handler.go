package schedule

import (
	"context"

	v1GenAPI "mrb-service/internal/generated/v1"
)

type ScheduleHandler interface {
	PostRoomsRoomIdScheduleCreate(ctx context.Context, request v1GenAPI.PostRoomsRoomIdScheduleCreateRequestObject) (v1GenAPI.PostRoomsRoomIdScheduleCreateResponseObject, error)
}

type scheduleHandler struct{}

func NewHandler() *scheduleHandler {
	return &scheduleHandler{}
}

func (h *scheduleHandler) PostRoomsRoomIdScheduleCreate(
	ctx context.Context,
	request v1GenAPI.PostRoomsRoomIdScheduleCreateRequestObject,
) (v1GenAPI.PostRoomsRoomIdScheduleCreateResponseObject, error) {
	return nil, nil
}
