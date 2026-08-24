package testutil

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"mrb-service/internal/config"
)

func NewTestPostgres(ctx context.Context, cfg config.Database, url string) (*pgxpool.Pool, error) {
	poolConfig, err := pgxpool.ParseConfig(url)
	if err != nil {
		return nil, fmt.Errorf(
			"parse database URL: %w",
			err,
		)
	}

	poolConfig.MaxConns = cfg.GetMaxOpenConns()
	poolConfig.MinConns = cfg.GetMinOpenConns()
	poolConfig.MaxConnLifetime = cfg.GetConnMaxLifetime()
	poolConfig.MaxConnIdleTime = cfg.GetMaxConnIdleTime()

	connectCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	pool, err := pgxpool.NewWithConfig(
		connectCtx,
		poolConfig,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"create PostgreSQL pool: %w",
			err,
		)
	}

	if err := pool.Ping(connectCtx); err != nil {
		pool.Close()

		return nil, fmt.Errorf(
			"ping PostgreSQL: %w",
			err,
		)
	}

	return pool, nil
}
