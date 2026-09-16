# Protobuf contracts

Source of truth for gRPC APIs and shared event envelope types.

## Layout

```text
api/proto/
├── buf.yaml
├── buf.gen.yaml
├── SUBJECTS.md              # NATS JetStream subjects
└── helpdesk/
    ├── common/v1/events.proto
    └── ticket/v1/ticket.proto

api/gen/go/                  # generated Go stubs (committed)
└── helpdesk/...
```

## Generate

Requires `buf`, `protoc-gen-go`, `protoc-gen-go-grpc` on `PATH`:

```bash
make proto
# or
cd api/proto && buf generate
```

## Subjects

See [SUBJECTS.md](SUBJECTS.md). Constants also live in `libs/eventkit/subjects`.
