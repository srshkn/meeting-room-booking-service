package middleware

import (
	"context"
	"errors"
	"log/slog"
	"net/http"

	v1GenAPI "mrb-service/internal/generated/v1"
)

var (
	ErrUnauthorized = errors.New("unauthorized")
	ErrForbidden    = errors.New("forbidden")
)

var operationPermissions = map[string]string{

	// -------------------------------------------------------------------------
	// Room

	"PostRoomsCreate": "room:create",
	"GetRoomsList":    "room:list",

	// -------------------------------------------------------------------------
	// Schedule

	"PostRoomsRoomIdScheduleCreate": "schedule:create",

	// -------------------------------------------------------------------------
	// Slot

	"GetRoomsRoomIdSlotsList": "slot:list",

	// -------------------------------------------------------------------------
	// Booking

	"PostBookingsCreate":          "booking:create",
	"GetBookingsList":             "booking:list",
	"GetBookingsMy":               "booking:my",
	"PostBookingsBookingIdCancel": "booking:cancel",
}

func Authorization(logger *slog.Logger) v1GenAPI.StrictMiddlewareFunc {
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
			requiredPermission, protected := operationPermissions[operationID]
			if !protected {
				return next(ctx, w, r, request)
			}

			principal, ok := PrincipalFromContext(ctx)
			if !ok {
				logger.Warn(
					"principal missing from context",
					slog.String("method", r.Method),
					slog.String("path", r.URL.Path),
					slog.String("remote_addr", r.RemoteAddr),
				)

				return nil, NewHTTPError(
					http.StatusUnauthorized,
					v1GenAPI.UNAUTHORIZED,
					"authentication required",
				)
			}

			for _, permission := range principal.Permissions {
				if permission == requiredPermission {
					return next(ctx, w, r, request)
				}
			}

			logger.Warn(
				"access denied",
				slog.String("operation", operationID),
				slog.String("user_id", principal.UserID),
				slog.String("role", principal.Role),
				slog.String("required_permission", requiredPermission),
			)

			return nil, NewHTTPError(
				http.StatusForbidden,
				v1GenAPI.FORBIDDEN,
				"insufficient permissions",
			)
		}
	}
}
