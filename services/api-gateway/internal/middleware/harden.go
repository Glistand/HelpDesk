package middleware

import (
	"context"
	"log/slog"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"google.golang.org/grpc/metadata"
)

const (
	HeaderCorrelationID = "X-Correlation-Id"
	HeaderRequestID     = "X-Request-Id"
)

type corrKey struct{}
type reqIDKey struct{}

// CorrelationIDFromContext returns the HTTP correlation id.
func CorrelationIDFromContext(ctx context.Context) string {
	v, _ := ctx.Value(corrKey{}).(string)
	return v
}

// Chain applies middleware outermost-first.
func Chain(h http.Handler, mws ...func(http.Handler) http.Handler) http.Handler {
	for i := len(mws) - 1; i >= 0; i-- {
		h = mws[i](h)
	}
	return h
}

// SecureHeaders sets basic hardening headers.
func SecureHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "DENY")
		w.Header().Set("Referrer-Policy", "no-referrer")
		next.ServeHTTP(w, r)
	})
}

// MaxBody limits request body size (default 1 MiB).
func MaxBody(n int64) func(http.Handler) http.Handler {
	if n <= 0 {
		n = 1 << 20
	}
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			r.Body = http.MaxBytesReader(w, r.Body, n)
			next.ServeHTTP(w, r)
		})
	}
}

// Timeout cancels the request context after d.
func Timeout(d time.Duration) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx, cancel := context.WithTimeout(r.Context(), d)
			defer cancel()
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// Correlation ensures X-Correlation-Id / X-Request-Id and forwards them to gRPC.
func Correlation(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cid := strings.TrimSpace(r.Header.Get(HeaderCorrelationID))
		if cid == "" {
			cid = uuid.NewString()
		}
		rid := strings.TrimSpace(r.Header.Get(HeaderRequestID))
		if rid == "" {
			rid = uuid.NewString()
		}
		w.Header().Set(HeaderCorrelationID, cid)
		w.Header().Set(HeaderRequestID, rid)
		ctx := context.WithValue(r.Context(), corrKey{}, cid)
		ctx = context.WithValue(ctx, reqIDKey{}, rid)
		ctx = metadata.AppendToOutgoingContext(ctx,
			"x-correlation-id", cid,
			"x-request-id", rid,
		)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// AccessLog logs method, path, status, duration, correlation id.
func AccessLog(logger *slog.Logger) func(http.Handler) http.Handler {
	if logger == nil {
		logger = slog.Default()
	}
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			rw := &statusWriter{ResponseWriter: w, status: 200}
			next.ServeHTTP(rw, r)
			logger.Info("http request",
				"method", r.Method,
				"path", r.URL.Path,
				"status", rw.status,
				"duration_ms", time.Since(start).Milliseconds(),
				"correlation_id", CorrelationIDFromContext(r.Context()),
			)
		})
	}
}

type statusWriter struct {
	http.ResponseWriter
	status int
}

func (w *statusWriter) WriteHeader(code int) {
	w.status = code
	w.ResponseWriter.WriteHeader(code)
}

// RateLimit is a simple per-IP token bucket (capacity=burst, refill=rate/sec).
func RateLimit(rate float64, burst int) func(http.Handler) http.Handler {
	type bucket struct {
		tokens float64
		last   time.Time
	}
	var mu sync.Mutex
	buckets := map[string]*bucket{}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ip := clientIP(r)
			now := time.Now()
			mu.Lock()
			b, ok := buckets[ip]
			if !ok {
				b = &bucket{tokens: float64(burst), last: now}
				buckets[ip] = b
			}
			elapsed := now.Sub(b.last).Seconds()
			b.tokens += elapsed * rate
			if b.tokens > float64(burst) {
				b.tokens = float64(burst)
			}
			b.last = now
			if b.tokens < 1 {
				mu.Unlock()
				w.Header().Set("Retry-After", "1")
				http.Error(w, `{"error":"rate limit exceeded"}`, http.StatusTooManyRequests)
				return
			}
			b.tokens--
			mu.Unlock()
			next.ServeHTTP(w, r)
		})
	}
}

func clientIP(r *http.Request) string {
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		parts := strings.Split(xff, ",")
		return strings.TrimSpace(parts[0])
	}
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return host
}

// WithOTel wraps the handler with OpenTelemetry HTTP instrumentation.
func WithOTel(service string, next http.Handler) http.Handler {
	return otelhttp.NewHandler(next, service,
		otelhttp.WithSpanNameFormatter(func(_ string, r *http.Request) string {
			return r.Method + " " + r.URL.Path
		}),
	)
}
