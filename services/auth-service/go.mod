module github.com/Glistand/HelpDesk/services/auth-service

go 1.26

require (
	github.com/Glistand/HelpDesk/api/gen/go v0.0.0
	github.com/Glistand/HelpDesk/libs/grpckit v0.0.0
	github.com/golang-jwt/jwt/v5 v5.2.1
	golang.org/x/crypto v0.55.0
	google.golang.org/grpc v1.83.2
)

require (
	github.com/google/uuid v1.6.0 // indirect
	golang.org/x/net v0.58.0 // indirect
	golang.org/x/sys v0.47.0 // indirect
	golang.org/x/text v0.41.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20260526163538-3dc84a4a5aaa // indirect
	google.golang.org/protobuf v1.36.12 // indirect
)

replace (
	github.com/Glistand/HelpDesk/api/gen/go => ../../api/gen/go
	github.com/Glistand/HelpDesk/libs/grpckit => ../../libs/grpckit
)
