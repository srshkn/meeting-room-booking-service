package schedule

import (
	"context"

	v1GenAPI "mrb-service/internal/generated/v1"
)

type ScheduleHandler struct{}

func (h *ScheduleHandler) PostRoomsRoomIdScheduleCreate(
	ctx context.Context,
	request v1GenAPI.PostRoomsRoomIdScheduleCreateRequestObject,
) (v1GenAPI.PostRoomsRoomIdScheduleCreateResponseObject, error) {
	return nil, nil
}
