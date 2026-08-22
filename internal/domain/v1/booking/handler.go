package booking

import (
	"context"

	v1GenAPI "mrb-service/internal/generated/v1"
)

type BookingHandler interface {
	PostBookingsCreate(ctx context.Context, request v1GenAPI.PostBookingsCreateRequestObject) (v1GenAPI.PostBookingsCreateResponseObject, error)
	GetBookingsList(ctx context.Context, request v1GenAPI.GetBookingsListRequestObject) (v1GenAPI.GetBookingsListResponseObject, error)
	GetBookingsMy(ctx context.Context, request v1GenAPI.GetBookingsMyRequestObject) (v1GenAPI.GetBookingsMyResponseObject, error)
	PostBookingsBookingIdCancel(ctx context.Context, request v1GenAPI.PostBookingsBookingIdCancelRequestObject) (v1GenAPI.PostBookingsBookingIdCancelResponseObject, error)
}

type bookingHandler struct{}

func NewHandler() *bookingHandler {
	return &bookingHandler{}
}

func (h *bookingHandler) PostBookingsCreate(
	ctx context.Context,
	request v1GenAPI.PostBookingsCreateRequestObject,
) (v1GenAPI.PostBookingsCreateResponseObject, error) {
	return nil, nil
}

func (h *bookingHandler) GetBookingsList(
	ctx context.Context,
	request v1GenAPI.GetBookingsListRequestObject,
) (v1GenAPI.GetBookingsListResponseObject, error) {
	return nil, nil
}

func (h *bookingHandler) GetBookingsMy(
	ctx context.Context,
	request v1GenAPI.GetBookingsMyRequestObject,
) (v1GenAPI.GetBookingsMyResponseObject, error) {
	return nil, nil
}

func (h *bookingHandler) PostBookingsBookingIdCancel(
	ctx context.Context,
	request v1GenAPI.PostBookingsBookingIdCancelRequestObject,
) (v1GenAPI.PostBookingsBookingIdCancelResponseObject, error) {
	return nil, nil
}
