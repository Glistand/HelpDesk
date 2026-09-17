package main

import (
	"context"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"syscall"

	authv1 "github.com/Glistand/HelpDesk/api/gen/go/helpdesk/auth/v1"
	"github.com/Glistand/HelpDesk/libs/grpckit"
	"github.com/Glistand/HelpDesk/libs/otelkit"
	"github.com/Glistand/HelpDesk/services/auth-service/internal/config"
	"github.com/Glistand/HelpDesk/services/auth-service/internal/grpcserver"
	"github.com/Glistand/HelpDesk/services/auth-service/internal/store"
	"github.com/Glistand/HelpDesk/services/auth-service/internal/tokens"
	"google.golang.org/grpc"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	cfg := config.Load()

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	otelShutdown, err := otelkit.Init(ctx, "auth-service")
	if err != nil {
		logger.Error("otel init failed", "error", err)
		os.Exit(1)
	}
	defer func() { _ = otelShutdown(context.Background()) }()

	users := store.NewWithSeed()
	issuer := tokens.NewIssuer(cfg.JWTSecret, cfg.TokenTTL)

	lis, err := net.Listen("tcp", cfg.GRPCAddr)
	if err != nil {
		logger.Error("listen failed", "error", err)
		os.Exit(1)
	}

	srv := grpc.NewServer(grpckit.DefaultServerOptions(logger)...)
	authv1.RegisterAuthServiceServer(srv, grpcserver.New(users, issuer))

	go func() {
		<-ctx.Done()
		srv.GracefulStop()
	}()

	logger.Info("auth-service listening", "addr", cfg.GRPCAddr)
	if err := srv.Serve(lis); err != nil {
		logger.Error("serve failed", "error", err)
		os.Exit(1)
	}
}
