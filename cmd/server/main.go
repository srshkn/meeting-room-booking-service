package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"mrb-service/internal/app"
	"mrb-service/internal/config"
	"mrb-service/internal/db"
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

	pool, err := postgres.NewPool(ctx, cfg.DB)
	if err != nil {
		logger.Error(
			"connect to PostgreSQL: %v",
			slog.Any("error", err),
		)
		os.Exit(1)
	}
	defer pool.Close()

	queries := db.New(pool)

	// -------------------------------------------------------------------------
	// Server

	serverApp := app.New(
		cfg.Server,
		logger,
		queries,
		cfg.CORS,
	)

	// -------------------------------------------------------------------------
	// Run

	err = serverApp.Run(ctx)
	if err != nil {
		logger.Error(
			"server stopped unexpectedly",
			slog.Any("error", err),
		)
	}
}
