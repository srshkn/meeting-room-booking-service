package booking

import (
	"context"

	v1GenAPI "mrb-service/internal/generated/v1"
)

type BookingHandler struct{}

func (h *BookingHandler) PostBookingsCreate(
	ctx context.Context,
	request v1GenAPI.PostBookingsCreateRequestObject,
) (v1GenAPI.PostBookingsCreateResponseObject, error)

func (h *BookingHandler) GetBookingsList(
	ctx context.Context,
	request v1GenAPI.GetBookingsListRequestObject,
) (v1GenAPI.GetBookingsListResponseObject, error)

func (h *BookingHandler) GetBookingsMy(
	ctx context.Context,
	request v1GenAPI.GetBookingsMyRequestObject,
) (v1GenAPI.GetBookingsMyResponseObject, error)

func (h *BookingHandler) PostBookingsBookingIdCancel(
	ctx context.Context,
	request v1GenAPI.PostBookingsBookingIdCancelRequestObject,
) (v1GenAPI.PostBookingsBookingIdCancelResponseObject, error)
