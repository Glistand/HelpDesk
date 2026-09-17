// Package otelkit initializes OpenTelemetry tracing for HelpDesk services.
package otelkit

import (
	"context"
	"os"
	"strings"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.37.0"
)

// Init configures a global TracerProvider that exports via OTLP/HTTP.
// If OTEL_EXPORTER_OTLP_ENDPOINT is empty or OTEL_SDK_DISABLED=true, Init is a no-op.
func Init(ctx context.Context, serviceName string) (shutdown func(context.Context) error, err error) {
	noop := func(context.Context) error { return nil }

	if os.Getenv("OTEL_SDK_DISABLED") == "true" {
		return noop, nil
	}
	endpoint := strings.TrimSpace(os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT"))
	if endpoint == "" {
		return noop, nil
	}
	if name := strings.TrimSpace(os.Getenv("OTEL_SERVICE_NAME")); name != "" {
		serviceName = name
	}

	opts := []otlptracehttp.Option{otlptracehttp.WithInsecure()}
	switch {
	case strings.HasPrefix(endpoint, "http://"), strings.HasPrefix(endpoint, "https://"):
		opts = append(opts, otlptracehttp.WithEndpointURL(strings.TrimRight(endpoint, "/")+"/v1/traces"))
	default:
		opts = append(opts, otlptracehttp.WithEndpoint(endpoint))
	}

	exp, err := otlptracehttp.New(ctx, opts...)
	if err != nil {
		return nil, err
	}

	res, err := resource.New(ctx,
		resource.WithFromEnv(),
		resource.WithTelemetrySDK(),
		resource.WithHost(),
		resource.WithAttributes(semconv.ServiceName(serviceName)),
	)
	if err != nil {
		return nil, err
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exp),
		sdktrace.WithResource(res),
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
	)
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))

	return tp.Shutdown, nil
}
