package slot

import (
	"context"

	v1GenAPI "mrb-service/internal/generated/v1"
)

type SlotHandler struct{}

func (h *SlotHandler) GetRoomsRoomIdSlotsList(
	ctx context.Context,
	request v1GenAPI.GetRoomsRoomIdSlotsListRequestObject,
) (v1GenAPI.GetRoomsRoomIdSlotsListResponseObject, error) {
	return nil, nil
}
