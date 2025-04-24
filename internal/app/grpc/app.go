package grpcapp

import (
	"fmt"
	"log/slog"
	"net"

	authgrpc "github.com/VariableSan/gia-sso/internal/grpc/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type App struct {
	log        *slog.Logger
	gRPCServer *grpc.Server
	port       int
}

func New(
	log *slog.Logger,
	port int,
) *App {
	gRPCServer := grpc.NewServer()

	authgrpc.Register(gRPCServer)
	reflection.Register(gRPCServer)

	return &App{
		log:        log,
		gRPCServer: gRPCServer,
		port:       port,
	}
}

func (app *App) MustRun() {
	if err := app.Run(); err != nil {
		panic(err)
	}
}

func (app *App) Run() error {
	const operation = "grpcapp.Run"

	log := app.log.With(
		slog.String("operation", operation),
		slog.Int("port", app.port),
	)

	listener, err := net.Listen("tcp", "localhost:"+fmt.Sprintf("%d", app.port))
	if err != nil {
		return fmt.Errorf("%s: %w", operation, err)
	}

	log.Info("grpc server is running", slog.String("addr", listener.Addr().String()))

	if err := app.gRPCServer.Serve(listener); err != nil {
		return fmt.Errorf("%s: %w", operation, err)
	}

	return nil
}

func (app *App) Stop() {
	const operation = "grpcapp.Stop"

	app.log.
		With(slog.String("operation", operation)).
		Info("stopping gRPC server", slog.Int("port", app.port))

	app.gRPCServer.GracefulStop()
}
