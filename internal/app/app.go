package app

import (
	"log/slog"
	"time"

	grpcapp "github.com/VariableSan/gia-sso/internal/app/grpc"
	"github.com/VariableSan/gia-sso/internal/services/auth"
	"github.com/VariableSan/gia-sso/internal/storage/sqlite"
)

type App struct {
	GRPCSrv *grpcapp.App
}

func New(
	log *slog.Logger,
	grpcHost string,
	grpcPort int,
	storagePath string,
	tokenTTL time.Duration,
) *App {
	storage, err := sqlite.New(storagePath)
	if err != nil {
		panic(err)
	}

	authService := auth.New(log, storage, tokenTTL)

	grpcApp := grpcapp.New(log, authService, grpcHost, grpcPort)

	return &App{
		GRPCSrv: grpcApp,
	}
}
