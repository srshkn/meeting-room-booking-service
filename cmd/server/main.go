package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"mrb-service/internal/config"
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

	cfg, help, err := config.New()
	if err != nil {
		fmt.Printf("help: %s;\nerr: %s\n", help, err)
		return
	}

	// -------------------------------------------------------------------------
	// Database

	pool, err := postgres.NewPool(ctx, cfg.Postgres)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer pool.Close()

	fmt.Printf("Server: host: %s, port: %d", cfg.Server.Host, cfg.Server.Port)
}
