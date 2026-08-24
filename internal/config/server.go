package config

import (
	"fmt"
	"time"
)

const (
	serverHostEnv            string = "SERVER_HOST"
	serverPortEnv            string = "SERVER_PORT"
	serverShutdownTimeoutEnv string = "SERVER_SHUTDOWN_TIMEOUT"
)

type Server interface {
	GetHost() string
	GetPort() int
	GetShutdownTimeout() time.Duration
}

type server struct {
	Host            string        `conf:"required"`
	Port            int           `conf:"required"`
	ShutdownTimeout time.Duration `conf:"required"`
}

func (s server) GetHost() string {
	return s.Host
}

func (s server) GetPort() int {
	return s.Port
}

func (s server) GetShutdownTimeout() time.Duration {
	return s.ShutdownTimeout
}

func (s server) validate() error {
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
			serverShutdownTimeoutEnv,
		)
	}

	return nil
}
