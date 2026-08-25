package auth

import (
	"context"
	"errors"

	"mrb-service/internal/cookie"

	v1GenAPI "mrb-service/internal/generated/v1"
)

type AuthHandler interface {
	PostAuthDummyLogin(ctx context.Context, request v1GenAPI.PostAuthDummyLoginRequestObject) (v1GenAPI.PostAuthDummyLoginResponseObject, error)
	PostAuthLogin(ctx context.Context, request v1GenAPI.PostAuthLoginRequestObject) (v1GenAPI.PostAuthLoginResponseObject, error)
}

type authHandler struct {
	service AuthService
	cookie  cookie.Auth
}

func NewHandler(service AuthService, cookie cookie.Auth) *authHandler {
	return &authHandler{
		service: service,
		cookie:  cookie,
	}
}

func (h *authHandler) PostAuthDummyLogin(
	ctx context.Context,
	request v1GenAPI.PostAuthDummyLoginRequestObject,
) (v1GenAPI.PostAuthDummyLoginResponseObject, error) {
	return nil, nil
}

func (h *authHandler) PostAuthLogin(
	ctx context.Context,
	request v1GenAPI.PostAuthLoginRequestObject,
) (v1GenAPI.PostAuthLoginResponseObject, error) {
	token, refreshToken, err := h.service.Login(ctx, *request.Body)
	if err != nil {

		switch {
		case errors.Is(err, ErrIncorrectPassword),
			errors.Is(err, ErrUserNotFound),
			errors.Is(err, ErrEmptyEmail),
			errors.Is(err, ErrLongEmail),
			errors.Is(err, ErrEmptyPassword),
			errors.Is(err, ErrShortPassword),
			errors.Is(err, ErrLongPassword):

			response := v1GenAPI.PostAuthLogin401JSONResponse{}
			response.Error.Code = v1GenAPI.UNAUTHORIZED
			response.Error.Message = err.Error()

			return response, nil

		default:
			response := v1GenAPI.PostAuthLogin500JSONResponse{}
			response.Error.Code = string(v1GenAPI.INTERNALERROR) // "INTERNAL_ERROR"
			response.Error.Message = "internal server error"

			return response, nil
		}

	}

	return v1GenAPI.PostAuthLogin200JSONResponse{
		Body: v1GenAPI.Token{
			Token: token.Token,
		},
		Headers: v1GenAPI.PostAuthLogin200ResponseHeaders{
			SetCookie: h.cookie.SetRefreshToken(refreshToken.Token, refreshToken.ExpiresAt),
		},
	}, nil
}
