package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Glistand/HelpDesk/services/api-gateway/internal/clients"
	"github.com/Glistand/HelpDesk/services/api-gateway/internal/config"
	"github.com/Glistand/HelpDesk/services/api-gateway/internal/handlers"
	"github.com/Glistand/HelpDesk/services/api-gateway/internal/middleware"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	cfg := config.Load()

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	c, err := clients.Dial(ctx, cfg.AuthGRPCAddr, cfg.TicketGRPCAddr, cfg.AssignmentGRPCAddr, cfg.AuditGRPCAddr)
	if err != nil {
		logger.Error("dial grpc failed", "error", err)
		os.Exit(1)
	}
	defer c.Close()

	api := handlers.New(c)
	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", api.Health)
	mux.HandleFunc("POST /auth/login", api.Login)

	auth := middleware.Auth(c.Auth)
	mux.Handle("POST /tickets", auth(http.HandlerFunc(api.CreateTicket)))
	mux.Handle("GET /tickets", auth(http.HandlerFunc(api.ListTickets)))
	mux.Handle("GET /tickets/{id}", auth(http.HandlerFunc(api.GetTicket)))
	mux.Handle("GET /tickets/{id}/timeline", auth(http.HandlerFunc(api.GetTimeline)))
	mux.Handle("PATCH /tickets/{id}/status", auth(http.HandlerFunc(api.UpdateTicketStatus)))

	srv := &http.Server{
		Addr:              cfg.HTTPAddr,
		Handler:           withCORS(mux),
		ReadHeaderTimeout: 5 * time.Second,
	}

	go func() {
		<-ctx.Done()
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = srv.Shutdown(shutdownCtx)
	}()

	logger.Info("api-gateway listening", "addr", cfg.HTTPAddr)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Error("serve failed", "error", err)
		os.Exit(1)
	}
}

func withCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PATCH, OPTIONS")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}
