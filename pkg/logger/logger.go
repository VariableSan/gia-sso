package logger

import (
	"log/slog"
	"os"

	"github.com/VariableSan/gia-sso/internal/config"
	"github.com/VariableSan/gia-sso/pkg/logger/prettylog"
)

func NewLogger(env string) *slog.Logger {
	var log *slog.Logger

	timeOnly := func(groups []string, a slog.Attr) slog.Attr {
		if a.Key == slog.TimeKey {
			a.Value = slog.StringValue(a.Value.Time().Format("15:04:05")) // Format time as "hh:mm:ss"
		}
		return a
	}

	switch env {
	case config.EnvLocal:
		log = slog.New(prettylog.NewHandler(nil))
	case config.EnvDev:
		log = slog.New(prettylog.NewHandler(nil))
	case config.EnvProd:
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
				Level:       slog.LevelInfo,
				ReplaceAttr: timeOnly,
			}),
		)
	}

	return log
}
