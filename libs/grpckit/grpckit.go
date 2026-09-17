package grpckit

import (
	"log/slog"

	"github.com/Glistand/HelpDesk/libs/grpckit/interceptors"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
)

// DefaultServerOptions returns StatsHandler (OTel) + recovery/correlation/logging.
func DefaultServerOptions(logger *slog.Logger) []grpc.ServerOption {
	return []grpc.ServerOption{
		grpc.StatsHandler(otelgrpc.NewServerHandler()),
		grpc.UnaryInterceptor(interceptors.UnaryServerChain(
			interceptors.Recovery(logger),
			interceptors.Correlation(),
			interceptors.Logging(logger),
		)),
	}
}

// DefaultUnaryServerInterceptors is kept for compatibility; prefer DefaultServerOptions.
func DefaultUnaryServerInterceptors(logger *slog.Logger) grpc.ServerOption {
	return grpc.UnaryInterceptor(interceptors.UnaryServerChain(
		interceptors.Recovery(logger),
		interceptors.Correlation(),
		interceptors.Logging(logger),
	))
}

// DefaultClientOptions returns OTel stats handler + correlation propagation.
func DefaultClientOptions() []grpc.DialOption {
	return []grpc.DialOption{
		grpc.WithStatsHandler(otelgrpc.NewClientHandler()),
		grpc.WithUnaryInterceptor(interceptors.UnaryClientPropagate()),
	}
}

// DefaultUnaryClientInterceptors is kept for compatibility; prefer DefaultClientOptions.
func DefaultUnaryClientInterceptors() grpc.DialOption {
	return grpc.WithUnaryInterceptor(interceptors.UnaryClientPropagate())
}
