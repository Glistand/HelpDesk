package main

import (
	"context"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	searchv1 "github.com/Glistand/HelpDesk/api/gen/go/helpdesk/search/v1"
	ticketv1 "github.com/Glistand/HelpDesk/api/gen/go/helpdesk/ticket/v1"
	"github.com/Glistand/HelpDesk/libs/eventkit/natsx"
	"github.com/Glistand/HelpDesk/libs/grpckit"
	"github.com/Glistand/HelpDesk/libs/otelkit"
	"github.com/Glistand/HelpDesk/services/search-service/internal/config"
	"github.com/Glistand/HelpDesk/services/search-service/internal/consumer"
	"github.com/Glistand/HelpDesk/services/search-service/internal/grpcserver"
	"github.com/Glistand/HelpDesk/services/search-service/internal/index"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	cfg := config.Load()

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	otelShutdown, err := otelkit.Init(ctx, "search-service")
	if err != nil {
		logger.Error("otel init failed", "error", err)
		os.Exit(1)
	}
	defer func() { _ = otelShutdown(context.Background()) }()

	idx, err := index.New(cfg.MeiliHost, cfg.MeiliAPIKey, cfg.MeiliIndex)
	if err != nil {
		logger.Error("meili client failed", "error", err)
		os.Exit(1)
	}
	if err := waitMeili(ctx, idx, logger); err != nil {
		logger.Error("meili ensure failed", "error", err)
		os.Exit(1)
	}

	os.Setenv("NATS_URL", cfg.NATSURL)
	nc, js, err := natsx.Connect(ctx)
	if err != nil {
		logger.Error("nats connect failed", "error", err)
		os.Exit(1)
	}
	defer nc.Close()

	clientOpts := append([]grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}, grpckit.DefaultClientOptions()...)
	ticketConn, err := grpc.NewClient(cfg.TicketGRPCAddr, clientOpts...)
	if err != nil {
		logger.Error("ticket dial failed", "error", err)
		os.Exit(1)
	}
	defer ticketConn.Close()

	h := consumer.New(ticketv1.NewTicketServiceClient(ticketConn), idx, logger)
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
	searchv1.RegisterSearchServiceServer(srv, grpcserver.New(idx))

	go func() {
		<-ctx.Done()
		srv.GracefulStop()
	}()

	logger.Info("search-service listening", "addr", cfg.GRPCAddr, "meili", cfg.MeiliHost)
	if err := srv.Serve(lis); err != nil {
		logger.Error("serve failed", "error", err)
		os.Exit(1)
	}
}

func waitMeili(ctx context.Context, idx *index.Store, logger *slog.Logger) error {
	var last error
	for i := 0; i < 30; i++ {
		if err := idx.Ensure(ctx); err != nil {
			last = err
			logger.Warn("meili not ready", "attempt", i+1, "error", err)
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(time.Second):
			}
			continue
		}
		return nil
	}
	return last
}
