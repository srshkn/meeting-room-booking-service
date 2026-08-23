package hash

import (
	"errors"
)

var (
	ErrInvalidHash        = errors.New("invalid password hash")
	ErrUnsupportedHash    = errors.New("unsupported password hash")
	ErrUnsupportedVersion = errors.New("unsupported argon2 version")
)
