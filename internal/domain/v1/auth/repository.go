package auth

import (
	"context"

	"github.com/google/uuid"

	"mrb-service/internal/db"
)

type repository interface {
	GetUserByLogin(ctx context.Context, email string) (db.GetUserByLoginRow, error)
	CreateRefreshToken(ctx context.Context, arg db.CreateRefreshTokenParams) (uuid.UUID, error)
	GetUserByID(ctx context.Context, id uuid.UUID) (db.GetUserByIDRow, error)
}
