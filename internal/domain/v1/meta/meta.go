package meta

import (
	"context"

	v1GenAPI "mrb-service/internal/generated/v1"
)

type MetaHandler struct{}

func New() *MetaHandler {
	return &MetaHandler{}
}

func (h *MetaHandler) GetHealth(
	ctx context.Context,
	request v1GenAPI.GetHealthRequestObject,
) (v1GenAPI.GetHealthResponseObject, error) {
	return v1GenAPI.GetHealth200JSONResponse{
		Status: "OK",
	}, nil
}
