package user

import (
	"context"

	"mrb-service/internal/db"
)

type repository interface {
	CreateUser(ctx context.Context, arg db.CreateUserParams) (db.CreateUserRow, error)
}
