package config

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
)

const (
	authCookieNameEnv     string = "AUTH_COOKIE_NAME"
	authCookiePathEnv     string = "AUTH_COOKIE_PATH"
	authCookieDomainEnv   string = "AUTH_COOKIE_DOMAIN"
	authCookieSecureEnv   string = "AUTH_COOKIE_SECURE"
	authCookieHttpOnlyEnv string = "AUTH_COOKIE_HTTP_ONLY"
	authCookieSameSiteEnv string = "AUTH_COOKIE_SAME_SITE"
	authCookieMaxAgeEnv   string = "AUTH_COOKIE_MAX_AGE"
)

type Cookie interface {
	GetRefreshTokenName() string
	GetPath() string
	GetDomain() string
	GetSecure() bool
	GetHTTPOnly() bool
	GetSameSite() http.SameSite
}

type cookie struct {
	RefreshTokenName string `conf:"required"`
	Path             string `conf:"required"`
	Domain           string `conf:"required"`
	Secure           bool   `conf:"required"`
	HTTPOnly         bool   `conf:"required"`
	StringSameSite   string `conf:"required"`
	sameSite         http.SameSite
}

func (c *cookie) GetRefreshTokenName() string {
	return c.RefreshTokenName
}

func (c *cookie) GetPath() string {
	return c.Path
}

func (c *cookie) GetDomain() string {
	return c.Domain
}

func (c *cookie) GetSecure() bool {
	return c.Secure
}

func (c *cookie) GetHTTPOnly() bool {
	return c.HTTPOnly
}

func (c *cookie) GetSameSite() http.SameSite {
	return c.sameSite
}

func (c *cookie) validate() error {

	switch {
	case c.RefreshTokenName == "":
		return errors.New("refresh token cookie name is empty")

	case !strings.HasPrefix(c.Path, "/"):
		return fmt.Errorf("invalid cookie path: %q", c.Path)

	case c.Domain != "":
		if strings.TrimSpace(c.Domain) != c.Domain {
			return errors.New("cookie domain contains leading or trailing spaces")
		}

		if strings.ContainsAny(c.Domain, "/ \t\r\n") {
			return errors.New("cookie domain contains invalid characters")
		}
	}

	switch c.sameSite {
	case http.SameSiteDefaultMode,
		http.SameSiteLaxMode,
		http.SameSiteStrictMode,
		http.SameSiteNoneMode:
	default:
		return errors.New("invalid cookie same site")
	}

	if c.sameSite == http.SameSiteNoneMode && !c.Secure {
		return errors.New("cookie secure must be enabled when same site is none")
	}

	return nil
}

func (c *cookie) parseSameSite() error {
	var err error

	switch strings.ToLower(c.StringSameSite) {
	case "default":
		c.sameSite = http.SameSiteDefaultMode

	case "lax":
		c.sameSite = http.SameSiteLaxMode

	case "strict":
		c.sameSite = http.SameSiteStrictMode

	case "none":
		c.sameSite = http.SameSiteNoneMode

	default:
		err = fmt.Errorf("invalid SameSite value %q", c.StringSameSite)
	}

	return err
}
