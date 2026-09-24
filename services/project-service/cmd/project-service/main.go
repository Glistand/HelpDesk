package main

import (
	"context"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"syscall"

	projectv1 "github.com/Glistand/HelpDesk/api/gen/go/helpdesk/project/v1"
	"github.com/Glistand/HelpDesk/libs/grpckit"
	"github.com/Glistand/HelpDesk/libs/otelkit"
	"github.com/Glistand/HelpDesk/services/project-service/internal/config"
	"github.com/Glistand/HelpDesk/services/project-service/internal/db"
	"github.com/Glistand/HelpDesk/services/project-service/internal/grpcserver"
	"github.com/Glistand/HelpDesk/services/project-service/internal/repository"
	"google.golang.org/grpc"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	cfg := config.Load()
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	shutdown, err := otelkit.Init(ctx, "project-service")
	if err != nil {
		logger.Error("otel init failed", "error", err)
		os.Exit(1)
	}
	defer func() { _ = shutdown(context.Background()) }()

	sqlDB, err := db.Open(ctx, cfg.DatabaseURL)
	if err != nil {
		logger.Error("db open failed", "error", err)
		os.Exit(1)
	}
	defer sqlDB.Close()
	if err := db.Migrate(ctx, sqlDB); err != nil {
		logger.Error("migrate failed", "error", err)
		os.Exit(1)
	}
	repo := repository.New(sqlDB)
	if err := repo.Bootstrap(ctx, cfg.Bootstrap); err != nil {
		logger.Error("bootstrap profile failed", "error", err)
		os.Exit(1)
	}

	listener, err := net.Listen("tcp", cfg.GRPCAddr)
	if err != nil {
		logger.Error("listen failed", "error", err)
		os.Exit(1)
	}
	server := grpc.NewServer(grpckit.DefaultServerOptions(logger)...)
	projectv1.RegisterProjectServiceServer(server, grpcserver.New(repo))
	go func() {
		<-ctx.Done()
		server.GracefulStop()
	}()
	logger.Info("project-service listening", "addr", cfg.GRPCAddr)
	if err := server.Serve(listener); err != nil {
		logger.Error("serve failed", "error", err)
		os.Exit(1)
	}
}
