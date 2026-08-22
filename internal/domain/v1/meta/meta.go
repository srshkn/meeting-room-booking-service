package meta

import (
	"context"

	v1GenAPI "mrb-service/internal/generated/v1"
)

type MetaHandler interface {
	GetHealth(ctx context.Context, request v1GenAPI.GetHealthRequestObject) (v1GenAPI.GetHealthResponseObject, error)
}

type metaHandler struct{}

func NewHandler() *metaHandler {
	return &metaHandler{}
}

func (h *metaHandler) GetHealth(
	ctx context.Context,
	request v1GenAPI.GetHealthRequestObject,
) (v1GenAPI.GetHealthResponseObject, error) {
	return v1GenAPI.GetHealth200JSONResponse{
		Status: "OK",
	}, nil
}
