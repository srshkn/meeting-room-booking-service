package middleware

import (
	"context"
	"log/slog"
	v1GenAPI "mrb-service/internal/generated/v1"
	"net/http"
	"time"
)

func Logging(logger *slog.Logger) v1GenAPI.StrictMiddlewareFunc {
	return func(
		next v1GenAPI.StrictHandlerFunc,
		operationID string,
	) v1GenAPI.StrictHandlerFunc {
		return func(
			ctx context.Context,
			w http.ResponseWriter,
			r *http.Request,
			request any,
		) (any, error) {
			startedAt := time.Now()

			response, err := next(ctx, w, r, request)

			logger.Info(
				"HTTP request",
				slog.String("operation", operationID),
				slog.String("method", r.Method),
				slog.String("path", r.URL.Path),
				slog.Duration("duration", time.Since(startedAt)),
				slog.Any("error", err),
				slog.String("remote_addr", r.RemoteAddr),
			)

			return response, err
		}
	}
}
