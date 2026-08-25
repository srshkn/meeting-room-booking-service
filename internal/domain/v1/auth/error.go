package auth

import (
	"errors"
)

var (
	ErrIncorrectPassword = errors.New("incorrect password")
	ErrUserNotFound      = errors.New("user not found")

	ErrEmptyEmail = errors.New("empty email")
	ErrLongEmail  = errors.New("long email")

	ErrEmptyPassword = errors.New("empty password")
	ErrShortPassword = errors.New("short password")
	ErrLongPassword  = errors.New("long password")
)
