package config

import (
	"fmt"
	"net/url"
	"time"
)

const (
	postgresHostEnv     string = "POSTGRES_HOST"
	postgresPortEnv     string = "POSTGRES_PORT"
	postgresUserEnv     string = "POSTGRES_USER"
	postgresPasswordEnv string = "POSTGRES_PASSWORD"
	postgresDataBaseEnv string = "POSTGRES_DB"

	maxOpenConnsEnv    string = "MAX_OPEN_CONNS"
	minOpenConnsEnv    string = "MIN_OPEN_CONNS"
	connMaxLifetimeEnv string = "CONN_MAX_LIFETIME"
	maxConnIdleTimeEnv string = "MAX_CONN_IDLE_TIME"
)

type Database interface {
	GetURL() string
	GetMaxOpenConns() int32
	GetMinOpenConns() int32
	GetConnMaxLifetime() time.Duration
	GetMaxConnIdleTime() time.Duration
}

type Postgres struct {
	Host     string `conf:"required"`
	Port     int    `conf:"required"`
	User     string `conf:"required"`
	Password string `conf:"required,mask"`
	Database string `conf:"required"`

	MaxOpenConns    int32         `conf:"required"`
	MinOpenConns    int32         `conf:"required"`
	ConnMaxLifetime time.Duration `conf:"required"`
	MaxConnIdleTime time.Duration `conf:"required"`
}

func (p Postgres) GetURL() string {
	databaseURL := url.URL{
		Scheme: "postgres",
		User:   url.UserPassword(p.User, p.Password),
		Host:   fmt.Sprintf("%s:%d", p.Host, p.Port),
		Path:   p.Database,
	}

	query := databaseURL.Query()
	query.Set("sslmode", "disable")
	databaseURL.RawQuery = query.Encode()

	return databaseURL.String()
}

func (p Postgres) GetMaxOpenConns() int32 {
	return p.MaxOpenConns
}

func (p Postgres) GetMinOpenConns() int32 {
	return p.MinOpenConns
}

func (p Postgres) GetConnMaxLifetime() time.Duration {
	return p.ConnMaxLifetime
}

func (p Postgres) GetMaxConnIdleTime() time.Duration {
	return p.MaxConnIdleTime
}

func (p Postgres) validatePostgres() error {
	switch {

	case p.Host == "":
		return fmt.Errorf("environment variable %q is required", postgresHostEnv)

	case p.Port < 1 || p.Port > 65535:
		return fmt.Errorf("environment variable %q must be between 1 and 65535", postgresPortEnv)

	case p.User == "":
		return fmt.Errorf("environment variable %q is required", postgresUserEnv)

	case p.Password == "":
		return fmt.Errorf("environment variable %q is required", postgresPasswordEnv)

	case p.Database == "":
		return fmt.Errorf("environment variable %q is required", postgresDataBaseEnv)

	case p.MaxOpenConns <= 0:
		return fmt.Errorf(
			"environment variable %q must be greater than 0",
			maxOpenConnsEnv,
		)

	case p.MinOpenConns <= 0:
		return fmt.Errorf(
			"environment variable %q must be greater than 0",
			minOpenConnsEnv,
		)

	case p.ConnMaxLifetime <= 0:
		return fmt.Errorf(
			"environment variable %q must be greater than 0",
			connMaxLifetimeEnv,
		)

	case p.MaxConnIdleTime <= 0:
		return fmt.Errorf(
			"environment variable %q must be greater than 0",
			maxConnIdleTimeEnv,
		)

	case p.MinOpenConns > p.MaxOpenConns:
		return fmt.Errorf(
			"environment variable %q must not exceed %q",
			minOpenConnsEnv,
			maxOpenConnsEnv,
		)

	}

	if len(p.User) > 63 {
		return fmt.Errorf(
			"environment variable %q must not exceed 63 bytes",
			postgresUserEnv,
		)
	}

	if len(p.Database) > 63 {
		return fmt.Errorf(
			"environment variable %q must not exceed 63 bytes",
			postgresDataBaseEnv,
		)
	}

	return nil
}
