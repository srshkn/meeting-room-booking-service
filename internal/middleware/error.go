package middleware

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"

	v1GenAPI "mrb-service/internal/generated/v1"
)

type HTTPError struct {
	Status int
	Body   v1GenAPI.ErrorResponse
}

func (e *HTTPError) Error() string {
	return e.Body.Error.Message
}

func NewHTTPError(
	status int,
	code v1GenAPI.ErrorResponseErrorCode,
	message string,
) error {
	response := v1GenAPI.ErrorResponse{}
	response.Error.Code = code
	response.Error.Message = message

	return &HTTPError{
		Status: status,
		Body:   response,
	}
}

func WriteResponseError(
	w http.ResponseWriter,
	r *http.Request,
	err error,
) {
	var httpErr *HTTPError

	if !errors.As(err, &httpErr) {
		httpErr = &HTTPError{
			Status: http.StatusInternalServerError,
		}

		httpErr.Body.Error.Code = v1GenAPI.INTERNALERROR
		httpErr.Body.Error.Message = "internal server error"
	}

	var buf bytes.Buffer
	if encodeErr := json.NewEncoder(&buf).Encode(httpErr.Body); encodeErr != nil {
		http.Error(
			w,
			"internal server error",
			http.StatusInternalServerError,
		)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpErr.Status)
	_, _ = buf.WriteTo(w)
}
