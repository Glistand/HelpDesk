package natsx

import (
	"context"

	"github.com/Glistand/HelpDesk/libs/eventkit/envelope"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

func startConsumeSpan(ctx context.Context, ev envelope.Event, subject string) (context.Context, trace.Span) {
	tracer := otel.Tracer("helpdesk/eventkit")
	return tracer.Start(ctx, "nats.consume "+ev.Type,
		trace.WithSpanKind(trace.SpanKindConsumer),
		trace.WithAttributes(
			attribute.String("messaging.system", "nats"),
			attribute.String("messaging.destination", subject),
			attribute.String("messaging.operation", "process"),
			attribute.String("helpdesk.event_id", ev.EventID),
			attribute.String("helpdesk.correlation_id", ev.CorrelationID),
			attribute.String("helpdesk.aggregate_id", ev.AggregateID),
			attribute.String("helpdesk.event_type", ev.Type),
		),
	)
}

func endConsumeSpan(span trace.Span, err error) {
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
	}
	span.End()
}
