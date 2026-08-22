package user

import (
	"context"
	"errors"

	v1GenAPI "mrb-service/internal/generated/v1"
)

type UserHandler interface {
	PostUserRegister(ctx context.Context, request v1GenAPI.PostUserRegisterRequestObject) (v1GenAPI.PostUserRegisterResponseObject, error)
}

type userHandler struct {
	service UserService
}

func NewHandler(service UserService) *userHandler {
	return &userHandler{
		service: service,
	}
}

func (h *userHandler) PostUserRegister(
	ctx context.Context,
	request v1GenAPI.PostUserRegisterRequestObject,
) (v1GenAPI.PostUserRegisterResponseObject, error) {
	user, err := h.service.Registration(ctx, *request.Body)
	if err != nil {
		switch {

		case errors.Is(err, ErrEmailAlreadyExists), errors.Is(err, ErrUsernameAlreadyExists):
			response := v1GenAPI.PostUserRegister400JSONResponse{}
			response.Error.Code = v1GenAPI.INVALIDREQUEST
			response.Error.Message = "invalid request"

			return response, nil

		default:
			response := v1GenAPI.PostUserRegister500JSONResponse{}
			response.Error.Code = "INTERNAL_ERROR"
			response.Error.Message = "internal server error"

			return response, nil

		}
	}

	return v1GenAPI.PostUserRegister201JSONResponse{
		User: user,
	}, nil
}
