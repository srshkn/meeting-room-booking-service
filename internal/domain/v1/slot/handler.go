package slot

import (
	"context"

	v1GenAPI "mrb-service/internal/generated/v1"
)

type SlotHandler interface {
	GetRoomsRoomIdSlotsList(ctx context.Context, request v1GenAPI.GetRoomsRoomIdSlotsListRequestObject) (v1GenAPI.GetRoomsRoomIdSlotsListResponseObject, error)
}

type slotHandler struct{}

func NewHandler() *slotHandler {
	return &slotHandler{}
}

func (h *slotHandler) GetRoomsRoomIdSlotsList(
	ctx context.Context,
	request v1GenAPI.GetRoomsRoomIdSlotsListRequestObject,
) (v1GenAPI.GetRoomsRoomIdSlotsListResponseObject, error) {
	return nil, nil
}
