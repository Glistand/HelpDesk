package main

import (
	"context"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"syscall"

	slav1 "github.com/Glistand/HelpDesk/api/gen/go/helpdesk/sla/v1"
	"github.com/Glistand/HelpDesk/libs/eventkit/natsx"
	"github.com/Glistand/HelpDesk/libs/grpckit"
	"github.com/Glistand/HelpDesk/libs/otelkit"
	"github.com/Glistand/HelpDesk/services/sla-service/internal/config"
	"github.com/Glistand/HelpDesk/services/sla-service/internal/consumer"
	"github.com/Glistand/HelpDesk/services/sla-service/internal/db"
	"github.com/Glistand/HelpDesk/services/sla-service/internal/grpcserver"
	"github.com/Glistand/HelpDesk/services/sla-service/internal/policy"
	"github.com/Glistand/HelpDesk/services/sla-service/internal/repository"
	"github.com/Glistand/HelpDesk/services/sla-service/internal/timers"
	"github.com/redis/go-redis/v9"
	"google.golang.org/grpc"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	cfg := config.Load()

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	otelShutdown, err := otelkit.Init(ctx, "sla-service")
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

	opt, err := redis.ParseURL(cfg.RedisURL)
	if err != nil {
		logger.Error("redis url invalid", "error", err)
		os.Exit(1)
	}
	rdb := redis.NewClient(opt)
	if err := rdb.Ping(ctx).Err(); err != nil {
		logger.Error("redis ping failed", "error", err)
		os.Exit(1)
	}
	defer rdb.Close()

	os.Setenv("NATS_URL", cfg.NATSURL)
	nc, js, err := natsx.Connect(ctx)
	if err != nil {
		logger.Error("nats connect failed", "error", err)
		os.Exit(1)
	}
	defer nc.Close()

	repo := repository.New(sqlDB)
	store := timers.New(rdb, logger)
	pub := natsx.NewPublisher(js)
	pol := policy.Policy{
		Name:          "default",
		FirstResponse: cfg.FirstResponse,
		Resolve:       cfg.Resolve,
		WarnRatio:     cfg.WarnRatio,
	}

	h := consumer.New(repo, store, pub, pol, logger)
	if err := consumer.Subscribe(ctx, js, h); err != nil {
		logger.Error("subscribe failed", "error", err)
		os.Exit(1)
	}

	worker := timers.NewWorker(store, repo, pub, cfg.PollInterval, logger)
	go worker.Run(ctx)

	lis, err := net.Listen("tcp", cfg.GRPCAddr)
	if err != nil {
		logger.Error("listen failed", "error", err)
		os.Exit(1)
	}
	srv := grpc.NewServer(grpckit.DefaultServerOptions(logger)...)
	slav1.RegisterSLAServiceServer(srv, grpcserver.New(repo))

	go func() {
		<-ctx.Done()
		srv.GracefulStop()
	}()

	logger.Info("sla-service listening",
		"addr", cfg.GRPCAddr,
		"first_response", cfg.FirstResponse.String(),
		"resolve", cfg.Resolve.String(),
		"warn_ratio", cfg.WarnRatio,
	)
	if err := srv.Serve(lis); err != nil {
		logger.Error("serve failed", "error", err)
		os.Exit(1)
	}
}
