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
		case errors.Is(err, ErrEmptyUsername),
			errors.Is(err, ErrLongUsername),
			errors.Is(err, ErrEmptyEmail),
			errors.Is(err, ErrLongEmail),
			errors.Is(err, ErrEmptyPassword),
			errors.Is(err, ErrShortPassword),
			errors.Is(err, ErrLongPassword),
			errors.Is(err, ErrPasswordsDoNotMatch),
			errors.Is(err, ErrUsernameAlreadyExists),
			errors.Is(err, ErrEmailAlreadyExists):

			response := v1GenAPI.PostUserRegister400JSONResponse{}
			response.Error.Code = v1GenAPI.INVALIDREQUEST
			response.Error.Message = err.Error()

			return response, nil

		default:
			response := v1GenAPI.PostUserRegister500JSONResponse{}
			response.Error.Code = string(v1GenAPI.INTERNALERROR) // "INTERNAL_ERROR"
			response.Error.Message = "internal server error"

			return response, nil

		}
	}

	return v1GenAPI.PostUserRegister201JSONResponse{
		User: user,
	}, nil
}
