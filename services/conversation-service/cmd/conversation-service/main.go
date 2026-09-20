package main

import (
	"context"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"syscall"

	conversationv1 "github.com/Glistand/HelpDesk/api/gen/go/helpdesk/conversation/v1"
	"github.com/Glistand/HelpDesk/libs/eventkit/natsx"
	"github.com/Glistand/HelpDesk/libs/grpckit"
	"github.com/Glistand/HelpDesk/libs/otelkit"
	"github.com/Glistand/HelpDesk/services/conversation-service/internal/config"
	"github.com/Glistand/HelpDesk/services/conversation-service/internal/db"
	"github.com/Glistand/HelpDesk/services/conversation-service/internal/grpcserver"
	"github.com/Glistand/HelpDesk/services/conversation-service/internal/openrouter"
	"github.com/Glistand/HelpDesk/services/conversation-service/internal/repository"
	"google.golang.org/grpc"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	cfg := config.Load()

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	otelShutdown, err := otelkit.Init(ctx, "conversation-service")
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

	bot := &openrouter.Client{
		APIKey:  cfg.OpenRouterAPIKey,
		Model:   cfg.OpenRouterModel,
		BaseURL: cfg.OpenRouterBaseURL,
	}
	if bot.APIKey == "" {
		logger.Warn("OPENROUTER_API_KEY empty — bot replies will use fallback text")
	}

	repo := repository.New(sqlDB)
	pub := natsx.NewPublisher(js)
	srvImpl := grpcserver.New(repo, bot, pub, logger)

	lis, err := net.Listen("tcp", cfg.GRPCAddr)
	if err != nil {
		logger.Error("listen failed", "error", err)
		os.Exit(1)
	}

	srv := grpc.NewServer(grpckit.DefaultServerOptions(logger)...)
	conversationv1.RegisterConversationServiceServer(srv, srvImpl)

	go func() {
		<-ctx.Done()
		srv.GracefulStop()
	}()

	logger.Info("conversation-service listening", "addr", cfg.GRPCAddr)
	if err := srv.Serve(lis); err != nil {
		logger.Error("serve failed", "error", err)
		os.Exit(1)
	}
}
