package user

import (
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
)

const (
	uniqueViolationCode    string = "23505"
	usersUsernameUniqueIDX string = "users_username_unique_idx"
	usersEmailUniqueIDX    string = "users_email_unique_idx"
)

var (
	ErrEmptyUsername = errors.New("empty username")
	ErrLongUsername  = errors.New("long username")

	ErrEmptyEmail = errors.New("empty email")
	ErrLongEmail  = errors.New("long email")

	ErrEmptyPassword       = errors.New("empty password")
	ErrShortPassword       = errors.New("short password")
	ErrLongPassword        = errors.New("long password")
	ErrPasswordsDoNotMatch = errors.New("passwords do not match")

	ErrUsernameAlreadyExists = errors.New("it username already exists")
	ErrEmailAlreadyExists    = errors.New("it email already exists")

	ErrInvalidRequest = errors.New("invalid request")
)

func mapCreateUserError(err error) error {
	var pgErr *pgconn.PgError

	if !errors.As(err, &pgErr) {
		return err
	}

	if pgErr.Code != uniqueViolationCode {
		return err
	}

	switch pgErr.ConstraintName {
	case usersUsernameUniqueIDX:
		return ErrUsernameAlreadyExists

	case usersEmailUniqueIDX:
		return ErrEmailAlreadyExists

	default:
		return err
	}
}
