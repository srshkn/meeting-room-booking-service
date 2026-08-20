package room

import (
	"context"

	v1GenAPI "mrb-service/internal/generated/v1"
)

type RoomHandler struct{}

func (h *RoomHandler) PostRoomsCreate(
	ctx context.Context,
	request v1GenAPI.PostRoomsCreateRequestObject,
) (v1GenAPI.PostRoomsCreateResponseObject, error) {
	return nil, nil
}

func (h *RoomHandler) GetRoomsList(
	ctx context.Context,
	request v1GenAPI.GetRoomsListRequestObject,
) (v1GenAPI.GetRoomsListResponseObject, error) {
	return nil, nil
}
