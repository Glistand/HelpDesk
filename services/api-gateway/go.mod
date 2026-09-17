module github.com/Glistand/HelpDesk/services/api-gateway

go 1.26

require (
	github.com/Glistand/HelpDesk/api/gen/go v0.0.0-00010101000000-000000000000
	github.com/Glistand/HelpDesk/libs/grpckit v0.0.0
	github.com/Glistand/HelpDesk/libs/otelkit v0.0.0
	github.com/google/uuid v1.6.0
	go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp v0.64.0
	google.golang.org/grpc v1.83.2
)

require (
	github.com/cenkalti/backoff/v5 v5.0.3 // indirect
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/felixge/httpsnoop v1.0.4 // indirect
	github.com/go-logr/logr v1.4.4 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.27.3 // indirect
	go.opentelemetry.io/auto/sdk v1.2.1 // indirect
	go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc v0.64.0 // indirect
	go.opentelemetry.io/otel v1.44.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlptrace v1.39.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp v1.39.0 // indirect
	go.opentelemetry.io/otel/metric v1.44.0 // indirect
	go.opentelemetry.io/otel/sdk v1.44.0 // indirect
	go.opentelemetry.io/otel/trace v1.44.0 // indirect
	go.opentelemetry.io/proto/otlp v1.9.0 // indirect
	golang.org/x/net v0.58.0 // indirect
	golang.org/x/sys v0.47.0 // indirect
	golang.org/x/text v0.41.0 // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20260526163538-3dc84a4a5aaa // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20260526163538-3dc84a4a5aaa // indirect
	google.golang.org/protobuf v1.36.12 // indirect
)

replace (
	github.com/Glistand/HelpDesk/api/gen/go => ../../api/gen/go
	github.com/Glistand/HelpDesk/libs/grpckit => ../../libs/grpckit
	github.com/Glistand/HelpDesk/libs/otelkit => ../../libs/otelkit
)
