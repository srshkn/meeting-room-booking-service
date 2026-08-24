package middleware

import (
	"net/http"
	"slices"
	"strings"

	"mrb-service/internal/config"
	v1GenAPI "mrb-service/internal/generated/v1"
)

func CORS(cfg config.CORS) v1GenAPI.MiddlewareFunc {
	return func(
		next http.Handler,
	) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			origin := r.Header.Get("Origin")

			if slices.Contains(cfg.GetAllowedOrigins(), origin) {
				w.Header().Set("Access-Control-Allow-Origin", origin)
				w.Header().Set("Access-Control-Allow-Credentials", cfg.GetallowCredentials())
			}

			w.Header().Set(
				"Access-Control-Allow-Methods",
				strings.Join(cfg.GetAllowedMethods(), ","),
			)

			w.Header().Set(
				"Access-Control-Allow-Headers",
				strings.Join(cfg.GetAllowedHeaders(), ","),
			)

			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusNoContent)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
