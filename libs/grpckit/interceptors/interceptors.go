package interceptors

import (
	"context"
	"log/slog"
	"runtime/debug"
	"time"

	"github.com/Glistand/HelpDesk/libs/grpckit/metadata"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// UnaryServerChain combines unary server interceptors (outermost first).
func UnaryServerChain(interceptors ...grpc.UnaryServerInterceptor) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (any, error) {
		h := handler
		for i := len(interceptors) - 1; i >= 0; i-- {
			interceptor := interceptors[i]
			next := h
			h = func(ctx context.Context, req any) (any, error) {
				return interceptor(ctx, req, info, next)
			}
		}
		return h(ctx, req)
	}
}

// Recovery recovers from panics and returns Internal.
func Recovery(logger *slog.Logger) grpc.UnaryServerInterceptor {
	if logger == nil {
		logger = slog.Default()
	}
	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (resp any, err error) {
		defer func() {
			if r := recover(); r != nil {
				logger.Error("grpc panic recovered",
					"method", info.FullMethod,
					"panic", r,
					"stack", string(debug.Stack()),
					"correlation_id", metadata.CorrelationIDFromContext(ctx),
				)
				err = status.Errorf(codes.Internal, "internal server error")
			}
		}()
		return handler(ctx, req)
	}
}

// Correlation ensures correlation/request ids exist on every request.
func Correlation() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (any, error) {
		return handler(metadata.EnsureIDs(ctx), req)
	}
}

// Logging logs method, duration, code, and correlation id.
func Logging(logger *slog.Logger) grpc.UnaryServerInterceptor {
	if logger == nil {
		logger = slog.Default()
	}
	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (any, error) {
		start := time.Now()
		resp, err := handler(ctx, req)
		code := codes.OK
		if err != nil {
			code = status.Code(err)
		}
		logger.Info("grpc unary",
			"method", info.FullMethod,
			"code", code.String(),
			"duration_ms", time.Since(start).Milliseconds(),
			"correlation_id", metadata.CorrelationIDFromContext(ctx),
			"request_id", metadata.RequestIDFromContext(ctx),
		)
		return resp, err
	}
}

// UnaryClientPropagate forwards correlation/request ids to the server.
func UnaryClientPropagate() grpc.UnaryClientInterceptor {
	return func(
		ctx context.Context,
		method string,
		req, reply any,
		cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker,
		opts ...grpc.CallOption,
	) error {
		return invoker(metadata.AppendOutgoing(ctx), method, req, reply, cc, opts...)
	}
}
