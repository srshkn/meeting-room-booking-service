package config

import (
	"errors"
	"fmt"
	"net/url"
	"strings"
)

const (
	corsAllowedOriginsEnv     string = "CORS_ALLOWED_ORIGINS"
	corsAllowedMethodsEnv     string = "CORS_ALLOWED_METHODS"
	corsAllowedHeadersEnv     string = "CORS_ALLOWED_HEADERS"
	corsAllowedCredentialsEnv string = "CORS_ALLOW_CREDENTIALS"
)

type CORS interface {
	GetAllowedOrigins() []string
	GetAllowedMethods() []string
	GetAllowedHeaders() []string
	GetallowCredentials() string
}

type cors struct {
	AllowedOrigins   []string `conf:"required"`
	AllowedMethods   []string `conf:"required"`
	AllowedHeaders   []string `conf:"required"`
	AllowCredentials string   `conf:"required"`
}

func (c cors) GetAllowedOrigins() []string {
	return c.AllowedOrigins
}

func (c cors) GetAllowedMethods() []string {
	return c.AllowedMethods
}

func (c cors) GetAllowedHeaders() []string {
	return c.AllowedHeaders
}

func (c cors) GetallowCredentials() string {
	return c.AllowCredentials
}

func (c cors) validate() error {
	for _, origin := range c.AllowedOrigins {
		if origin == "" {
			return errors.New("CORS origin cannot be empty")
		}

		u, err := url.Parse(origin)
		if err != nil {
			return fmt.Errorf("invalid CORS origin %q: %w", origin, err)
		}

		if u.Scheme != "http" && u.Scheme != "https" {
			return fmt.Errorf("invalid CORS scheme in %q", origin)
		}

		if u.Hostname() == "" {
			return fmt.Errorf(
				"invalid CORS origin %q: host must be localhost",
				origin,
			)
		}
	}

	for _, method := range c.AllowedMethods {
		method = strings.TrimSpace(method)

		switch method {
		case "GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS":
			// OK
		default:
			return fmt.Errorf("invalid CORS method %q", method)
		}
	}

	for _, header := range c.AllowedHeaders {
		header = strings.TrimSpace(header)

		if header == "" {
			return errors.New("CORS header cannot be empty")
		}
	}

	switch c.AllowCredentials {
	case "true", "false":
	// OK
	default:
		return fmt.Errorf("invalid CORS credentials %s", c.AllowCredentials)
	}

	return nil
}
