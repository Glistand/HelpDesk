# grpckit

Shared gRPC helpers for Helpdesk services.

## Contents

- `interceptors` — recovery, logging, correlation id
- `metadata` — `authorization`, `x-correlation-id`, `x-request-id`
- `statuserr` — typed gRPC status helpers

## Usage

```go
import (
    "log/slog"

    "github.com/Glistand/HelpDesk/libs/grpckit"
    "google.golang.org/grpc"
)

srv := grpc.NewServer(grpckit.DefaultUnaryServerInterceptors(slog.Default()))
```
