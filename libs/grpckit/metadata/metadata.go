package metadata

import (
	"context"

	"github.com/google/uuid"
	"google.golang.org/grpc/metadata"
)

const (
	Authorization = "authorization"
	CorrelationID = "x-correlation-id"
	RequestID     = "x-request-id"
)

// CorrelationIDFromContext returns x-correlation-id from incoming metadata.
func CorrelationIDFromContext(ctx context.Context) string {
	return first(ctx, CorrelationID)
}

// RequestIDFromContext returns x-request-id from incoming metadata.
func RequestIDFromContext(ctx context.Context) string {
	return first(ctx, RequestID)
}

// EnsureIDs returns a context with correlation and request ids present.
// Missing values are generated as UUIDs.
func EnsureIDs(ctx context.Context) context.Context {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		md = metadata.MD{}
	} else {
		md = md.Copy()
	}

	if len(md.Get(CorrelationID)) == 0 {
		md.Set(CorrelationID, uuid.NewString())
	}
	if len(md.Get(RequestID)) == 0 {
		md.Set(RequestID, uuid.NewString())
	}

	return metadata.NewIncomingContext(ctx, md)
}

// AppendOutgoing copies correlation/request ids onto outgoing metadata.
func AppendOutgoing(ctx context.Context) context.Context {
	pairs := make([]string, 0, 4)
	if v := CorrelationIDFromContext(ctx); v != "" {
		pairs = append(pairs, CorrelationID, v)
	}
	if v := RequestIDFromContext(ctx); v != "" {
		pairs = append(pairs, RequestID, v)
	}
	if len(pairs) == 0 {
		return ctx
	}
	return metadata.AppendToOutgoingContext(ctx, pairs...)
}

func first(ctx context.Context, key string) string {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return ""
	}
	vals := md.Get(key)
	if len(vals) == 0 {
		return ""
	}
	return vals[0]
}
