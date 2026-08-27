package room

import (
	"context"
	"errors"

	v1GenAPI "mrb-service/internal/generated/v1"
)

type RoomHandler interface {
	GetRoomsList(ctx context.Context, request v1GenAPI.GetRoomsListRequestObject) (v1GenAPI.GetRoomsListResponseObject, error)
	PostRoomsCreate(ctx context.Context, request v1GenAPI.PostRoomsCreateRequestObject) (v1GenAPI.PostRoomsCreateResponseObject, error)
}

type roomHandler struct {
	svc RoomService
}

func NewHandler(svc RoomService) *roomHandler {
	return &roomHandler{
		svc: svc,
	}
}

func (h *roomHandler) PostRoomsCreate(
	ctx context.Context,
	request v1GenAPI.PostRoomsCreateRequestObject,
) (v1GenAPI.PostRoomsCreateResponseObject, error) {
	room, err := h.svc.CreateRoom(ctx, *request.Body)
	if err != nil {
		switch {
		case errors.Is(err, ErrInvalidRoomName):

			response := v1GenAPI.PostRoomsCreate400JSONResponse{}
			response.Error.Code = v1GenAPI.INVALIDREQUEST
			response.Error.Message = err.Error()

			return response, nil

		default:
			response := v1GenAPI.PostRoomsCreate500JSONResponse{}
			response.Error.Code = string(v1GenAPI.INTERNALERROR) // "INTERNAL_ERROR"
			response.Error.Message = "internal server error"

			return response, nil
		}
	}
	return v1GenAPI.PostRoomsCreate201JSONResponse{
		Room: room,
	}, nil
}

func (h *roomHandler) GetRoomsList(
	ctx context.Context,
	request v1GenAPI.GetRoomsListRequestObject,
) (v1GenAPI.GetRoomsListResponseObject, error) {
	rooms, err := h.svc.GetListRooms(ctx)
	if err != nil {
		return nil, err
	}

	return v1GenAPI.GetRoomsList200JSONResponse{
		Rooms: &rooms,
	}, nil
}
