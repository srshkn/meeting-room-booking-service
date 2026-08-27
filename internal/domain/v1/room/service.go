package room

import (
	"context"
	"math"

	"mrb-service/internal/db"

	v1GenAPI "mrb-service/internal/generated/v1"
)

type RoomService interface {
	CreateRoom(ctx context.Context, body v1GenAPI.PostRoomsCreateJSONRequestBody) (*v1GenAPI.Room, error)
	GetListRooms(ctx context.Context) ([]v1GenAPI.Room, error)
}

type roomService struct {
	repo repository
}

func NewService(repo repository) *roomService {
	return &roomService{
		repo: repo,
	}
}

func (s *roomService) validateCreateRoom(name string, cap *int) error {
	if name == "" {
		return ErrInvalidRoomName
	}

	if cap != nil &&
		(*cap < 0 || *cap > math.MaxInt32) {
		return ErrInvalidRoomCapacity
	}

	return nil
}

func (s *roomService) CreateRoom(
	ctx context.Context,
	body v1GenAPI.PostRoomsCreateJSONRequestBody,
) (*v1GenAPI.Room, error) {
	name := body.Name

	if err := s.validateCreateRoom(name, body.Capacity); err != nil {
		return nil, err
	}

	createRoom, err := s.repo.CreateRoom(
		ctx,
		db.CreateRoomParams{
			Name:        name,
			Description: body.Description,
			Capacity:    body.Capacity,
		},
	)
	if err != nil {
		return nil, err
	}

	return &v1GenAPI.Room{
		Id:          createRoom.ID,
		Name:        createRoom.Name,
		Description: createRoom.Description,
		Capacity:    createRoom.Capacity,
		CreatedAt:   &createRoom.CreatedAt,
	}, nil
}

func (s *roomService) GetListRooms(
	ctx context.Context,
) ([]v1GenAPI.Room, error) {
	listRooms := make([]v1GenAPI.Room, 0)

	rooms, err := s.repo.GetRoomList(ctx)
	if err != nil {
		return listRooms, err
	}

	for _, val := range rooms {
		listRooms = append(listRooms, v1GenAPI.Room{
			Id:          val.ID,
			Name:        val.Name,
			Description: val.Description,
			Capacity:    val.Capacity,
			CreatedAt:   &val.CreatedAt,
		})
	}

	return listRooms, nil
}
