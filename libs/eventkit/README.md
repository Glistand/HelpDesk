# eventkit

Shared NATS JetStream helpers for Helpdesk services.

## Contents

- `envelope` — event envelope (`event_id`, `correlation_id`, `occurred_at`, …)
- `subjects` — JetStream subject / stream name constants
- `natsx` — connect, publish, durable subscribe skeleton
- `health` — HTTP health check helpers

## Usage

```go
ctx := context.Background()
nc, js, err := natsx.Connect(ctx)
if err != nil { ... }
defer nc.Close()

pub := natsx.NewPublisher(js)
ev, _ := envelope.New("ticket.created", ticketID, correlationID, payload)
_ = pub.Publish(ctx, subjects.TicketCreated, ev)
```
