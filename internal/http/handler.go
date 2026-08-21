package http

import (
	"mrb-service/internal/domain/v1/auth"
	"mrb-service/internal/domain/v1/booking"
	"mrb-service/internal/domain/v1/meta"
	"mrb-service/internal/domain/v1/room"
	"mrb-service/internal/domain/v1/schedule"
	"mrb-service/internal/domain/v1/slot"
	"mrb-service/internal/domain/v1/user"
	v1GenAPI "mrb-service/internal/generated/v1"
)

type Handler struct {
	*meta.MetaHandler
	*user.UserHandler
	*auth.AuthHandler
	*room.RoomHandler
	*schedule.ScheduleHandler
	*slot.SlotHandler
	*booking.BookingHandler
}

func New(
	metaHandler *meta.MetaHandler,
	userHandler *user.UserHandler,
	authHandler *auth.AuthHandler,
	roomHandler *room.RoomHandler,
	scheduleHandler *schedule.ScheduleHandler,
	slotHandler *slot.SlotHandler,
	bookingHandler *booking.BookingHandler,
) *Handler {
	return &Handler{
		MetaHandler:     metaHandler,
		UserHandler:     userHandler,
		AuthHandler:     authHandler,
		RoomHandler:     roomHandler,
		ScheduleHandler: scheduleHandler,
		SlotHandler:     slotHandler,
		BookingHandler:  bookingHandler,
	}
}

var _ v1GenAPI.StrictServerInterface = (*Handler)(nil)
