module github.com/Glistand/HelpDesk/services/escalation-service

go 1.26

require (
	github.com/Glistand/HelpDesk/api/gen/go v0.0.0
	github.com/Glistand/HelpDesk/libs/eventkit v0.0.0
	github.com/Glistand/HelpDesk/libs/grpckit v0.0.0
	github.com/jackc/pgx/v5 v5.7.2
	github.com/nats-io/nats.go v1.39.1
	google.golang.org/grpc v1.83.2
)

require (
	github.com/google/uuid v1.6.0 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20240606120523-5a60cdf6a761 // indirect
	github.com/jackc/puddle/v2 v2.2.2 // indirect
	github.com/klauspost/compress v1.17.9 // indirect
	github.com/nats-io/nkeys v0.4.9 // indirect
	github.com/nats-io/nuid v1.0.1 // indirect
	golang.org/x/crypto v0.55.0 // indirect
	golang.org/x/net v0.58.0 // indirect
	golang.org/x/sync v0.22.0 // indirect
	golang.org/x/sys v0.47.0 // indirect
	golang.org/x/text v0.41.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20260526163538-3dc84a4a5aaa // indirect
	google.golang.org/protobuf v1.36.12 // indirect
)

replace (
	github.com/Glistand/HelpDesk/api/gen/go => ../../api/gen/go
	github.com/Glistand/HelpDesk/libs/eventkit => ../../libs/eventkit
	github.com/Glistand/HelpDesk/libs/grpckit => ../../libs/grpckit
)
