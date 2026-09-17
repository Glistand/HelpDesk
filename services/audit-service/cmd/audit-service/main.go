package main

import (
	"context"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"syscall"

	auditv1 "github.com/Glistand/HelpDesk/api/gen/go/helpdesk/audit/v1"
	"github.com/Glistand/HelpDesk/libs/eventkit/natsx"
	"github.com/Glistand/HelpDesk/libs/grpckit"
	"github.com/Glistand/HelpDesk/libs/otelkit"
	"github.com/Glistand/HelpDesk/services/audit-service/internal/config"
	"github.com/Glistand/HelpDesk/services/audit-service/internal/consumer"
	"github.com/Glistand/HelpDesk/services/audit-service/internal/db"
	"github.com/Glistand/HelpDesk/services/audit-service/internal/grpcserver"
	"github.com/Glistand/HelpDesk/services/audit-service/internal/repository"
	"google.golang.org/grpc"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	cfg := config.Load()

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	otelShutdown, err := otelkit.Init(ctx, "audit-service")
	if err != nil {
		logger.Error("otel init failed", "error", err)
		os.Exit(1)
	}
	defer func() { _ = otelShutdown(context.Background()) }()

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

	os.Setenv("NATS_URL", cfg.NATSURL)
	nc, js, err := natsx.Connect(ctx)
	if err != nil {
		logger.Error("nats connect failed", "error", err)
		os.Exit(1)
	}
	defer nc.Close()

	repo := repository.New(sqlDB)
	h := consumer.New(repo, logger)
	if err := consumer.Subscribe(ctx, js, h); err != nil {
		logger.Error("subscribe failed", "error", err)
		os.Exit(1)
	}

	lis, err := net.Listen("tcp", cfg.GRPCAddr)
	if err != nil {
		logger.Error("listen failed", "error", err)
		os.Exit(1)
	}

	srv := grpc.NewServer(grpckit.DefaultServerOptions(logger)...)
	auditv1.RegisterAuditServiceServer(srv, grpcserver.New(repo))

	go func() {
		<-ctx.Done()
		srv.GracefulStop()
	}()

	logger.Info("audit-service listening", "addr", cfg.GRPCAddr)
	if err := srv.Serve(lis); err != nil {
		logger.Error("serve failed", "error", err)
		os.Exit(1)
	}
}
