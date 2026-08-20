package user

import (
	"context"

	v1GenAPI "mrb-service/internal/generated/v1"
)

type UserHandler struct{}

func (h *UserHandler) PostUserRegister(
	ctx context.Context,
	request v1GenAPI.PostUserRegisterRequestObject,
) (v1GenAPI.PostUserRegisterResponseObject, error)
