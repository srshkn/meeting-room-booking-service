package config

import (
	"errors"
	"fmt"
	"os"

	"github.com/ardanlabs/conf/v3"
	"github.com/joho/godotenv"
)

const (
	configPrefixEnv = "CONFIG_PREFIX"
)

type config struct {
	Server *server
	DB     *db
	Logger *logger
	CORS   *cors
	JWT    *jwt
	Cookie *cookie
}

func New() (*config, string, error) {
	var cfg config

	if err := godotenv.Load(); err != nil && !errors.Is(err, os.ErrNotExist) {
		return nil, "", fmt.Errorf("load .env: %w", err)
	}

	prefix := os.Getenv(configPrefixEnv)

	help, err := conf.Parse(prefix, &cfg)
	if err != nil {
		return nil, help, fmt.Errorf("parse config: %w", err)
	}

	if err := cfg.Server.validate(); err != nil {
		return nil, "", fmt.Errorf("validate config Server: %w", err)
	}

	if err := cfg.Logger.validate(); err != nil {
		return nil, "", fmt.Errorf("validate config LogManager: %w", err)
	}

	if err := cfg.DB.validate(); err != nil {
		return nil, "", fmt.Errorf("validate config Postgers: %w", err)
	}

	if err := cfg.JWT.loadKeys(); err != nil {
		return nil, "", fmt.Errorf("load keys JWT: %w", err)
	}

	if err := cfg.JWT.validate(); err != nil {
		return nil, "", fmt.Errorf("validate config JWT: %w", err)
	}

	if err := cfg.Cookie.parseSameSite(); err != nil {
		return nil, "", fmt.Errorf("parse SameSite cookie: %w", err)
	}

	if err := cfg.Cookie.validate(); err != nil {
		return nil, "", fmt.Errorf("validate config Cookie: %w", err)
	}

	if err := cfg.CORS.validate(); err != nil {
		return nil, "", fmt.Errorf("validate config CORS: %w", err)
	}

	return &cfg, "", nil
}
