package config

import (
	"fmt"
	"net/url"
	"time"
)

const (
	dbHostEnv     string = "DB_HOST"
	dbPortEnv     string = "DB_PORT"
	dbUserEnv     string = "DB_USER"
	dbPasswordEnv string = "DB_PASSWORD"
	dbDataBaseEnv string = "DB_NAME"

	dbMaxOpenConnsEnv    string = "DB_MAX_OPEN_CONNS"
	dbMinOpenConnsEnv    string = "DB_MIN_OPEN_CONNS"
	dbConnMaxLifetimeEnv string = "DB_CONN_MAX_LIFETIME"
	dbMaxConnIdleTimeEnv string = "DB_MAX_CONN_IDLE_TIME"
)

type Database interface {
	GetURL() string
	GetMaxOpenConns() int32
	GetMinOpenConns() int32
	GetConnMaxLifetime() time.Duration
	GetMaxConnIdleTime() time.Duration
}

type DB struct {
	Host     string `conf:"required"`
	Port     int    `conf:"required"`
	User     string `conf:"required"`
	Password string `conf:"required,mask"`
	Name     string `conf:"required"`

	MaxOpenConns    int32         `conf:"required"`
	MinOpenConns    int32         `conf:"required"`
	ConnMaxLifetime time.Duration `conf:"required"`
	MaxConnIdleTime time.Duration `conf:"required"`
}

func (p DB) GetURL() string {
	databaseURL := url.URL{
		Scheme: "postgres",
		User:   url.UserPassword(p.User, p.Password),
		Host:   fmt.Sprintf("%s:%d", p.Host, p.Port),
		Path:   p.Name,
	}

	query := databaseURL.Query()
	query.Set("sslmode", "disable")
	databaseURL.RawQuery = query.Encode()

	return databaseURL.String()
}

func (p DB) GetMaxOpenConns() int32 {
	return p.MaxOpenConns
}

func (p DB) GetMinOpenConns() int32 {
	return p.MinOpenConns
}

func (p DB) GetConnMaxLifetime() time.Duration {
	return p.ConnMaxLifetime
}

func (p DB) GetMaxConnIdleTime() time.Duration {
	return p.MaxConnIdleTime
}

func (p DB) validate() error {
	switch {

	case p.Host == "":
		return fmt.Errorf("environment variable %q is required", dbHostEnv)

	case p.Port < 1 || p.Port > 65535:
		return fmt.Errorf("environment variable %q must be between 1 and 65535", dbPortEnv)

	case p.User == "":
		return fmt.Errorf("environment variable %q is required", dbUserEnv)

	case p.Password == "":
		return fmt.Errorf("environment variable %q is required", dbPasswordEnv)

	case p.Name == "":
		return fmt.Errorf("environment variable %q is required", dbDataBaseEnv)

	case p.MaxOpenConns <= 0:
		return fmt.Errorf(
			"environment variable %q must be greater than 0",
			dbMaxOpenConnsEnv,
		)

	case p.MinOpenConns <= 0:
		return fmt.Errorf(
			"environment variable %q must be greater than 0",
			dbMinOpenConnsEnv,
		)

	case p.ConnMaxLifetime <= 0:
		return fmt.Errorf(
			"environment variable %q must be greater than 0",
			dbConnMaxLifetimeEnv,
		)

	case p.MaxConnIdleTime <= 0:
		return fmt.Errorf(
			"environment variable %q must be greater than 0",
			dbMaxConnIdleTimeEnv,
		)

	case p.MinOpenConns > p.MaxOpenConns:
		return fmt.Errorf(
			"environment variable %q must not exceed %q",
			dbMinOpenConnsEnv,
			dbMaxOpenConnsEnv,
		)

	}

	if len(p.User) > 63 {
		return fmt.Errorf(
			"environment variable %q must not exceed 63 bytes",
			dbUserEnv,
		)
	}

	if len(p.Name) > 63 {
		return fmt.Errorf(
			"environment variable %q must not exceed 63 bytes",
			dbDataBaseEnv,
		)
	}

	return nil
}
