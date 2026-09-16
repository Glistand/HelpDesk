package grpckit

import (
	"log/slog"

	"github.com/Glistand/HelpDesk/libs/grpckit/interceptors"
	"google.golang.org/grpc"
)

// DefaultUnaryServerInterceptors returns recovery + correlation + logging.
func DefaultUnaryServerInterceptors(logger *slog.Logger) grpc.ServerOption {
	return grpc.UnaryInterceptor(interceptors.UnaryServerChain(
		interceptors.Recovery(logger),
		interceptors.Correlation(),
		interceptors.Logging(logger),
	))
}

// DefaultUnaryClientInterceptors returns client-side id propagation.
func DefaultUnaryClientInterceptors() grpc.DialOption {
	return grpc.WithUnaryInterceptor(interceptors.UnaryClientPropagate())
}
