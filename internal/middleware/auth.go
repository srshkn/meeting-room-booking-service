package middleware

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"strings"

	"mrb-service/internal/jwt"

	v1GenAPI "mrb-service/internal/generated/v1"
)

const (
	prefix string = "Bearer "
)

var protectedOperations = map[string]struct{}{
	"GetRoomsList":                  {},
	"PostRoomsCreate":               {},
	"PostRoomsRoomIdScheduleCreate": {},
	"GetRoomsRoomIdSlotsList":       {},
	"PostBookingsCreate":            {},
	"GetBookingsList":               {},
	"GetBookingsMy":                 {},
	"PostBookingsBookingIdCancel":   {},
}

func extractBearerToken(rawToken string) (string, error) {

	if rawToken == "" {
		return "", errors.New("missing authorization header")
	}

	if !strings.HasPrefix(rawToken, prefix) {
		return "", errors.New("invalid authorization scheme")
	}

	token := strings.TrimSpace(strings.TrimPrefix(rawToken, prefix))
	if token == "" {
		return "", errors.New("missing bearer token")
	}

	return token, nil
}

func Auth(logger *slog.Logger, jwtManager jwt.Manager) v1GenAPI.StrictMiddlewareFunc {
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
			if _, protected := protectedOperations[operationID]; !protected {
				return next(ctx, w, r, request)
			}

			tokenString, err := extractBearerToken(r.Header.Get("Authorization"))
			if err != nil {
				logger.Warn(
					"access token rejected",
					slog.String("method", r.Method),
					slog.String("path", r.URL.Path),
					slog.Any("error", err),
					slog.String("remote_addr", r.RemoteAddr),
				)
				return nil, NewHTTPError(
					http.StatusUnauthorized,
					v1GenAPI.UNAUTHORIZED,
					"authentication required",
				)
			}

			claims, err := jwtManager.ValidateAccessToken(tokenString)
			if err != nil {
				logger.Warn(
					"access token rejected",
					slog.String("method", r.Method),
					slog.String("path", r.URL.Path),
					slog.Any("error", err),
					slog.String("remote_addr", r.RemoteAddr),
				)

				return nil, NewHTTPError(
					http.StatusUnauthorized,
					v1GenAPI.UNAUTHORIZED,
					"invalid or expired access token",
				)
			}

			principal := Principal{
				UserID:      claims.UserID,
				Role:        claims.Role,
				Permissions: append([]string(nil), claims.Permissions...),
			}

			newCtx := withPrincipal(ctx, principal)

			return next(
				newCtx,
				w,
				r.WithContext(newCtx),
				request,
			)
		}
	}
}
