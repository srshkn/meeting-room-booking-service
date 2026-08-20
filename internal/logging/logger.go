package logging

import (
	"log/slog"
	"os"

	"mrb-service/internal/config"
)

func New(cfg config.Logger) *slog.Logger {
	var handler slog.Handler

	switch cfg.GetFormat() {

	case "dev", "local":
		handler = slog.NewTextHandler(
			os.Stdout,
			&slog.HandlerOptions{
				Level:     slog.LevelDebug,
				AddSource: true,
			},
		)

	case "prod":
		handler = slog.NewJSONHandler(
			os.Stdout,
			&slog.HandlerOptions{
				Level: slog.LevelInfo,
			},
		)

	}

	return slog.New(handler)
}
