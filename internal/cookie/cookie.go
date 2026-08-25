package cookie

import (
	"net/http"
	"time"

	"mrb-service/internal/config"
)

type Auth interface {
	Name() string
	SetRefreshToken(token string, expiresAt time.Time) string
	ClearRefreshToken() string
}

type manager struct {
	refreshTokenName string
	path             string
	domain           string
	secure           bool
	httpOnly         bool
	sameSite         http.SameSite
}

func New(cfg config.Cookie) *manager {
	return &manager{
		refreshTokenName: cfg.GetRefreshTokenName(),
		path:             cfg.GetPath(),
		domain:           cfg.GetDomain(),
		secure:           cfg.GetSecure(),
		httpOnly:         cfg.GetHTTPOnly(),
		sameSite:         cfg.GetSameSite(),
	}
}

func (m *manager) Name() string {
	return m.refreshTokenName
}

func (m *manager) SetRefreshToken(
	token string,
	expiresAt time.Time,
) string {
	return (&http.Cookie{
		Name:     m.refreshTokenName,
		Value:    token,
		Expires:  expiresAt,
		Path:     m.path,
		Domain:   m.domain,
		Secure:   m.secure,
		HttpOnly: m.httpOnly,
		SameSite: m.sameSite,
	}).String()
}

func (m *manager) ClearRefreshToken() string {
	return (&http.Cookie{
		Name:     m.refreshTokenName,
		Value:    "",
		MaxAge:   -1,
		Path:     m.path,
		Domain:   m.domain,
		Secure:   m.secure,
		HttpOnly: m.httpOnly,
		SameSite: m.sameSite,
	}).String()
}
