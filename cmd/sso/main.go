package main

import (
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/VariableSan/gia-sso/internal/app"
	"github.com/VariableSan/gia-sso/internal/config"
	"github.com/VariableSan/gia-sso/pkg/prettylog"
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

	go application.GRPCSrv.MustRun()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	stopSignal := <-stop

	log.Info(
		"stopping application",
		slog.String("signal", stopSignal.String()),
	)

	application.GRPCSrv.Stop()

	log.Info("application stopped")
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
		log = slog.New(prettylog.NewHandler(nil))
	case envDev:
		log = slog.New(prettylog.NewHandler(nil))
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
