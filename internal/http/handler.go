package http

import (
	"mrb-service/internal/domain/v1/auth"
	"mrb-service/internal/domain/v1/booking"
	"mrb-service/internal/domain/v1/room"
	"mrb-service/internal/domain/v1/schedule"
	"mrb-service/internal/domain/v1/slot"
	"mrb-service/internal/domain/v1/user"
	v1GenAPI "mrb-service/internal/generated/v1"
)

type Handler struct {
	*user.UserHandler
	*auth.AuthHandler
	*room.RoomHandler
	*schedule.ScheduleHandler
	*slot.SlotHandler
	*booking.BookingHandler
}

func New() *Handler {
	return &Handler{}
}

var _ v1GenAPI.StrictServerInterface = (*Handler)(nil)
