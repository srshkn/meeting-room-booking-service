package room

import (
	"errors"
)

var (
	ErrInvalidRoomName     = errors.New("invalid room name")
	ErrInvalidRoomCapacity = errors.New("invalid room capacity")
)
