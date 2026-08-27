package app

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"strconv"
	"time"

	"mrb-service/internal/config"
	"mrb-service/internal/cookie"
	"mrb-service/internal/db"
	"mrb-service/internal/domain/v1/auth"
	"mrb-service/internal/domain/v1/booking"
	"mrb-service/internal/domain/v1/meta"
	"mrb-service/internal/domain/v1/room"
	"mrb-service/internal/domain/v1/schedule"
	"mrb-service/internal/domain/v1/slot"
	"mrb-service/internal/domain/v1/user"
	"mrb-service/internal/jwt"
	"mrb-service/internal/middleware"
	"mrb-service/internal/swagger"

	v1GenAPI "mrb-service/internal/generated/v1"
	handlerApp "mrb-service/internal/handler"
)

type ServerApp interface {
	Run(ctx context.Context) error
	Handler() http.Handler
}

type serverApp struct {
	httpServer      *http.Server
	logger          *slog.Logger
	shutdownTimeout time.Duration
}

func New(
	cfg config.Server,
	logger *slog.Logger,
	db db.Querier,
	jwtManager jwt.Manager,
	cookieManager cookie.Auth,
	cors config.CORS,
) *serverApp {
	mux := http.NewServeMux()

	userService := user.NewService(db)
	authService := auth.NewService(db, jwtManager)
	roomService := room.NewService(db)

	v1Handler := handlerApp.NewHandler(
		meta.NewHandler(),
		user.NewHandler(userService),
		auth.NewHandler(authService, cookieManager),
		room.NewHandler(roomService),
		schedule.NewHandler(),
		slot.NewHandler(),
		booking.NewHandler(),
	)

	strictHandler := v1GenAPI.NewStrictHandlerWithOptions(
		v1Handler,
		[]v1GenAPI.StrictMiddlewareFunc{
			middleware.Authorization(logger),
			middleware.Auth(logger, jwtManager),
			middleware.Logging(logger),
		},
		v1GenAPI.StrictHTTPServerOptions{
			ResponseErrorHandlerFunc: middleware.WriteResponseError,
		},
	)

	apiHandler := v1GenAPI.HandlerWithOptions(
		strictHandler,
		v1GenAPI.ChiServerOptions{
			BaseURL: "/api/v1",
		},
	)

	swagger.Register(mux)

	mux.Handle("/", middleware.CORS(cors)(apiHandler))

	address := net.JoinHostPort(
		cfg.GetHost(),
		strconv.Itoa(cfg.GetPort()),
	)

	return &serverApp{
		httpServer: &http.Server{
			Addr:              address,
			Handler:           mux,
			ReadHeaderTimeout: 5 * time.Second,
			ReadTimeout:       15 * time.Second,
			WriteTimeout:      15 * time.Second,
			IdleTimeout:       60 * time.Second,
		},
		logger:          logger,
		shutdownTimeout: cfg.GetShutdownTimeout(),
	}
}

func (s *serverApp) gracefulStop(serverErr <-chan error) error {
	shutdownCtx, cancel := context.WithTimeout(
		context.Background(),
		s.shutdownTimeout,
	)
	defer cancel()

	s.logger.Info(
		"shutting down HTTP server",
		slog.Duration("timeout", s.shutdownTimeout),
	)

	if err := s.httpServer.Shutdown(shutdownCtx); err != nil {
		s.logger.Error(
			"graceful shutdown failed",
			slog.Any("error", err),
		)

		if closeErr := s.httpServer.Close(); closeErr != nil {
			s.logger.Error(
				"force close HTTP server failed",
				slog.Any("error", closeErr),
			)
		}

		return fmt.Errorf("shutdown HTTP server: %w", err)
	}

	err := <-serverErr
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("HTTP server stopped: %w", err)
	}

	s.logger.Info("HTTP server stopped")

	return nil
}

func (s *serverApp) Run(ctx context.Context) error {
	serverErr := make(chan error, 1)

	go func() {
		serverErr <- s.httpServer.ListenAndServe()
	}()

	s.logger.Info("starting server",
		slog.String("address", s.httpServer.Addr),
	)

	s.logger.Info("server endpoints",
		slog.String("api", "http://"+s.httpServer.Addr),
		slog.String("swagger", "http://"+s.httpServer.Addr+"/docs/"),
	)

	select {

	case err := <-serverErr:
		if errors.Is(err, http.ErrServerClosed) {
			return nil
		}
		return err

	case <-ctx.Done():
		s.logger.Info("shutdown signal received")
		return s.gracefulStop(serverErr)

	}
}

func (s *serverApp) Handler() http.Handler {
	return s.httpServer.Handler
}
