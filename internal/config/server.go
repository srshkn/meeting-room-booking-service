package config

import (
	"fmt"
	"time"
)

const (
	serverHostEnv            string = "SERVER_HOST"
	serverPortEnv            string = "SERVER_PORT"
	shutdownTimeoutSecondEnv string = "SHUTDOWN_TIMEOUT_SECOND"
)

type Server interface {
	GetHost() string
	GetPort() int
	GetShutdownTimeout() time.Duration
}

type ServerHTTP struct {
	Host            string        `conf:"required"`
	Port            int           `conf:"required"`
	ShutdownTimeout time.Duration `conf:"required"`
}

func (s ServerHTTP) GetHost() string {
	return s.Host
}

func (s ServerHTTP) GetPort() int {
	return s.Port
}

func (s ServerHTTP) GetShutdownTimeout() time.Duration {
	return s.ShutdownTimeout
}

func (s ServerHTTP) validate() error {
	switch {

	case s.Host == "":
		return fmt.Errorf(
			"environment variable %q is required",
			serverHostEnv,
		)

	case s.Port < 1 || s.Port > 65535:
		return fmt.Errorf(
			"environment variable %q must be between 1 and 65535",
			serverPortEnv,
		)

	case s.ShutdownTimeout <= 0:
		return fmt.Errorf(
			"environment variable %q must be greater than 0",
			shutdownTimeoutSecondEnv,
		)
	}

	return nil
}
