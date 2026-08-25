package test

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"testing"

	"mrb-service/internal/app"
	"mrb-service/internal/config"
	"mrb-service/internal/cookie"
	"mrb-service/internal/db"
	"mrb-service/internal/jwt"
	"mrb-service/internal/logging"
	"mrb-service/internal/testutil"
)

var (
	testServerApp  app.ServerApp
	testLogger     *slog.Logger
	testDB         db.Querier
	testJWTManager jwt.JWTManager
	testCORS       config.CORS
	testHandler    http.Handler
)

func TestMain(m *testing.M) {

	// -------------------------------------------------------------------------
	// Context

	ctx := context.Background()

	// -------------------------------------------------------------------------
	// Test Env

	if err := testutil.NewTestEnv(); err != nil {
		return
	}

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

	testLogger = logging.New(cfg.Logger)

	// -------------------------------------------------------------------------
	// JWT

	testJWTManager = jwt.New(cfg.JWT)

	// -------------------------------------------------------------------------
	// Cookie manager

	cookieManager := cookie.New(cfg.Cookie)

	// -------------------------------------------------------------------------
	// CORSCfg

	testCORS = cfg.CORS

	// -------------------------------------------------------------------------
	// Container PostgreSQL

	url, cleanup, err := testutil.NewTestContainer(ctx)
	if err != nil {
		if cleanup != nil {
			cleanup()
		}

		slog.Error(
			"failed to setup test PostgreSQL",
			slog.Any("error", err),
		)
		os.Exit(1)
	}

	// -------------------------------------------------------------------------
	// PostgreSQL

	pool, err := testutil.NewTestPostgres(ctx, cfg.DB, url)
	if err != nil {
		if pool != nil {
			pool.Close()
		}

		cleanup()
		slog.Error(
			"failed to setup test PostgreSQL",
			slog.Any("error", err),
		)
		os.Exit(1)
	}

	// -------------------------------------------------------------------------
	// Migrations

	if err = testutil.NewTestMigrations(url); err != nil {
		pool.Close()
		cleanup()
		slog.Error(
			"failed to setup test PostgreSQL",
			slog.Any("error", err),
		)
		os.Exit(1)
	}

	testDB = db.New(pool)

	// -------------------------------------------------------------------------
	// ServerApp

	testServerApp = app.New(
		cfg.Server,
		testLogger,
		testDB,
		testJWTManager,
		cookieManager,
		testCORS,
	)

	// -------------------------------------------------------------------------
	// Tests

	code := m.Run()

	// -------------------------------------------------------------------------
	// Cleanup

	pool.Close()
	cleanup()

	os.Exit(code)
}
