package auth

import (
	"context"

	v1GenAPI "mrb-service/internal/generated/v1"
)

type AuthHandler interface {
	PostAuthDummyLogin(ctx context.Context, request v1GenAPI.PostAuthDummyLoginRequestObject) (v1GenAPI.PostAuthDummyLoginResponseObject, error)
	PostAuthLogin(ctx context.Context, request v1GenAPI.PostAuthLoginRequestObject) (v1GenAPI.PostAuthLoginResponseObject, error)
}

type authHandler struct{}

func NewHandler() *authHandler {
	return &authHandler{}
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
	return nil, nil
}
