# otelkit

OpenTelemetry bootstrap for HelpDesk services.

- Reads `OTEL_EXPORTER_OTLP_ENDPOINT` (e.g. `http://jaeger:4318`)
- Optional `OTEL_SERVICE_NAME`, `OTEL_SDK_DISABLED=true` to skip
- Sets global TracerProvider + W3C TraceContext propagator
