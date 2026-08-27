package room

import (
	"context"
	"mrb-service/internal/db"
)

type repository interface {
	CreateRoom(ctx context.Context, arg db.CreateRoomParams) (db.Room, error)
	GetRoomList(ctx context.Context) ([]db.Room, error)
}
