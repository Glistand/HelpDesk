package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Glistand/HelpDesk/libs/otelkit"
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

	otelShutdown, err := otelkit.Init(ctx, "api-gateway")
	if err != nil {
		logger.Error("otel init failed", "error", err)
		os.Exit(1)
	}
	defer func() { _ = otelShutdown(context.Background()) }()

	c, err := clients.Dial(ctx,
		cfg.AuthGRPCAddr,
		cfg.TicketGRPCAddr,
		cfg.AssignmentGRPCAddr,
		cfg.AuditGRPCAddr,
		cfg.SLAGRPCAddr,
		cfg.SearchGRPCAddr,
		cfg.ConversationGRPCAddr,
		cfg.ProjectGRPCAddr,
	)
	if err != nil {
		logger.Error("dial grpc failed", "error", err)
		os.Exit(1)
	}
	defer c.Close()

	api := handlers.New(c, cfg.WidgetSiteKey)
	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", api.Health)
	mux.HandleFunc("POST /auth/login", api.Login)

	// Public widget API (no JWT).
	mux.HandleFunc("POST /widget/session", api.WidgetSession)
	mux.HandleFunc("POST /widget/conversations/{id}/messages", api.WidgetPostMessage)
	mux.HandleFunc("GET /widget/conversations/{id}/messages", api.WidgetListMessages)
	mux.HandleFunc("POST /widget/conversations/{id}/handoff", api.WidgetHandoff)

	auth := middleware.Auth(c.Auth)
	mux.Handle("GET /project", auth(http.HandlerFunc(api.GetProject)))
	mux.Handle("PATCH /project", auth(http.HandlerFunc(api.UpdateProject)))
	mux.Handle("POST /tickets", auth(http.HandlerFunc(api.CreateTicket)))
	mux.Handle("GET /tickets", auth(http.HandlerFunc(api.ListTickets)))
	mux.Handle("GET /tickets/{id}", auth(http.HandlerFunc(api.GetTicket)))
	mux.Handle("GET /tickets/{id}/card", auth(http.HandlerFunc(api.GetTicketCard)))
	mux.Handle("GET /tickets/{id}/timeline", auth(http.HandlerFunc(api.GetTimeline)))
	mux.Handle("GET /tickets/{id}/sla", auth(http.HandlerFunc(api.GetSLA)))
	mux.Handle("PATCH /tickets/{id}/status", auth(http.HandlerFunc(api.UpdateTicketStatus)))
	mux.Handle("GET /search", auth(http.HandlerFunc(api.SearchTickets)))

	mux.Handle("GET /conversations", auth(http.HandlerFunc(api.ListConversations)))
	mux.Handle("GET /conversations/{id}", auth(http.HandlerFunc(api.GetConversation)))
	mux.Handle("POST /conversations/{id}/messages", auth(http.HandlerFunc(api.AgentPostMessage)))
	mux.Handle("POST /conversations/{id}/resolve", auth(http.HandlerFunc(api.ResolveConversation)))

	handler := middleware.Chain(mux,
		middleware.SecureHeaders,
		middleware.MaxBody(1<<20),
		middleware.RateLimit(100, 200),
		middleware.Timeout(60*time.Second),
		middleware.Correlation,
		middleware.AccessLog(logger),
		withCORS,
	)
	handler = middleware.WithOTel("api-gateway", handler)

	srv := &http.Server{
		Addr:              cfg.HTTPAddr,
		Handler:           handler,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       30 * time.Second,
		WriteTimeout:      90 * time.Second,
		IdleTimeout:       60 * time.Second,
		MaxHeaderBytes:    1 << 20,
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
		w.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type, X-Correlation-Id, X-Request-Id, X-Visitor-Id, X-Site-Key")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PATCH, OPTIONS")
		w.Header().Set("Access-Control-Expose-Headers", "X-Correlation-Id, X-Request-Id")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}
