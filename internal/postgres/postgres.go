package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"mrb-service/internal/config"
)

func NewPool(
	ctx context.Context,
	cfg config.Database,
) (*pgxpool.Pool, error) {
	config, err := pgxpool.ParseConfig(cfg.GetURL())
	if err != nil {
		return nil, fmt.Errorf("parse database URL: %w", err)
	}

	config.MaxConns = cfg.GetMaxOpenConns()
	config.MinConns = cfg.GetMinOpenConns()
	config.MaxConnLifetime = cfg.GetConnMaxLifetime()
	config.MaxConnIdleTime = cfg.GetMaxConnIdleTime()

	connectCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	pool, err := pgxpool.NewWithConfig(connectCtx, config)
	if err != nil {
		return nil, fmt.Errorf("create PostgreSQL pool: %w", err)
	}

	if err := pool.Ping(connectCtx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("ping PostgreSQL: %w", err)
	}

	return pool, nil
}
