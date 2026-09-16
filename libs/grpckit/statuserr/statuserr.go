package statuserr

import (
	"context"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// NotFound returns a gRPC NotFound status error.
func NotFound(msg string) error {
	return status.Error(codes.NotFound, msg)
}

// InvalidArgument returns a gRPC InvalidArgument status error.
func InvalidArgument(msg string) error {
	return status.Error(codes.InvalidArgument, msg)
}

// Internal returns a gRPC Internal status error.
func Internal(msg string) error {
	return status.Error(codes.Internal, msg)
}

// FromError maps a plain error to a gRPC status when possible.
func FromError(err error) error {
	if err == nil {
		return nil
	}
	if _, ok := status.FromError(err); ok {
		return err
	}
	if errors.Is(err, context.Canceled) {
		return status.Error(codes.Canceled, err.Error())
	}
	if errors.Is(err, context.DeadlineExceeded) {
		return status.Error(codes.DeadlineExceeded, err.Error())
	}
	return status.Errorf(codes.Internal, "%v", err)
}

// IsNotFound reports whether err is a NotFound status.
func IsNotFound(err error) bool {
	st, ok := status.FromError(err)
	return ok && st.Code() == codes.NotFound
}

// CodeOf returns the gRPC code for err, or Unknown.
func CodeOf(err error) codes.Code {
	if err == nil {
		return codes.OK
	}
	if st, ok := status.FromError(err); ok {
		return st.Code()
	}
	if errors.Is(err, context.Canceled) {
		return codes.Canceled
	}
	if errors.Is(err, context.DeadlineExceeded) {
		return codes.DeadlineExceeded
	}
	return codes.Unknown
}
