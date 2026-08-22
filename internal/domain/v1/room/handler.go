package room

import (
	"context"

	v1GenAPI "mrb-service/internal/generated/v1"
)

type RoomHandler interface {
	GetRoomsList(ctx context.Context, request v1GenAPI.GetRoomsListRequestObject) (v1GenAPI.GetRoomsListResponseObject, error)
	PostRoomsCreate(ctx context.Context, request v1GenAPI.PostRoomsCreateRequestObject) (v1GenAPI.PostRoomsCreateResponseObject, error)
}

type roomHandler struct{}

func NewHandler() *roomHandler {
	return &roomHandler{}
}

func (h *roomHandler) GetRoomsList(
	ctx context.Context,
	request v1GenAPI.GetRoomsListRequestObject,
) (v1GenAPI.GetRoomsListResponseObject, error) {
	return nil, nil
}

func (h *roomHandler) PostRoomsCreate(
	ctx context.Context,
	request v1GenAPI.PostRoomsCreateRequestObject,
) (v1GenAPI.PostRoomsCreateResponseObject, error) {
	return nil, nil
}
