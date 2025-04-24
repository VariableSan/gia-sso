package main

import (
	"log/slog"
	"os"

	"github.com/VariableSan/gia-sso/internal/app"
	"github.com/VariableSan/gia-sso/internal/config"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	cfg := config.MustLoad()

	log := setupLogger(cfg.Env)

	log.Info(
		"starting application on port: ",
		slog.Int("port", cfg.GRPC.Port),
	)

	application := app.New(
		log,
		cfg.GRPC.Port,
		cfg.StoragePath,
		cfg.TokenTTL,
	)

	application.GRPCSrv.MustRun()
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	timeOnly := func(groups []string, a slog.Attr) slog.Attr {
		if a.Key == slog.TimeKey {
			a.Value = slog.StringValue(a.Value.Time().Format("15:04:05")) // Format time as "hh:mm:ss"
		}
		return a
	}

	switch env {
	case envLocal:
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
				Level:       slog.LevelDebug,
				ReplaceAttr: timeOnly,
			}),
		)
	case envDev:
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
				Level:       slog.LevelDebug,
				ReplaceAttr: timeOnly,
			}),
		)
	case envProd:
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
				Level:       slog.LevelInfo,
				ReplaceAttr: timeOnly,
			}),
		)
	}

	return log
}
