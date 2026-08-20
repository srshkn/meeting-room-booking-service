package config

import (
	"fmt"
	"time"
)

const (
	serverHostEnv            string = "SERVER_HOST"
	serverPortEnv            string = "SERVER_PORT"
	shutdownTimeoutSecondEnv string = "SHOTDOWN_TIMEOUT_SECOND"
)

type iServer interface {
	GetHost() string
	GetPort() int
	GetShutdownTimeout() time.Duration
}

type Server struct {
	Host                  string        `conf:"required"`
	Port                  int           `conf:"required"`
	ShutdownTimeoutSecond time.Duration `conf:"required"`
}

func (s *Server) GetHost() string {
	return s.Host
}

func (s *Server) GetPort() int {
	return s.Port
}

func (s *Server) ShutdownTimeout() time.Duration {
	return s.ShutdownTimeoutSecond
}

func (s *Server) validateServer() error {
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

	case s.ShutdownTimeoutSecond <= 0:
		return fmt.Errorf(
			"environment variable %q must be greater than 0",
			shutdownTimeoutSecondEnv,
		)
	}

	return nil
}
