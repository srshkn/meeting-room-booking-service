package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"mrb-service/internal/config"
	"mrb-service/internal/logging"
	"mrb-service/internal/postgres"
)

func main() {

	// -------------------------------------------------------------------------
	// Context

	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stop()

	// -------------------------------------------------------------------------
	// Configuration

	cfg, _, err := config.New()
	if err != nil {
		slog.Error(
			"failed to load configuration",
			slog.Any("error", err),
		)
		return
	}

	// -------------------------------------------------------------------------
	// Logger

	logger := logging.New(cfg.Logger)

	// -------------------------------------------------------------------------
	// Database

	pool, err := postgres.NewPool(ctx, cfg.Postgres)
	if err != nil {
		logger.Error(
			"connect to PostgreSQL: %v",
			slog.Any("error", err),
		)
		os.Exit(1)
	}
	defer pool.Close()

	fmt.Printf("Server: host: %s, port: %d", cfg.Server.Host, cfg.Server.Port)
}
