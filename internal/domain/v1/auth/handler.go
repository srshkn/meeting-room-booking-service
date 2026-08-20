package auth

import (
	"context"

	v1GenAPI "mrb-service/internal/generated/v1"
)

type AuthHandler struct{}

func (h *AuthHandler) PostAuthDummyLogin(
	ctx context.Context,
	request v1GenAPI.PostAuthDummyLoginRequestObject,
) (v1GenAPI.PostAuthDummyLoginResponseObject, error)

func PostAuthLogin(
	ctx context.Context,
	request v1GenAPI.PostAuthLoginRequestObject,
) (v1GenAPI.PostAuthLoginResponseObject, error)
