module github.com/Glistand/HelpDesk/services/project-service

go 1.26

require (
	github.com/Glistand/HelpDesk/api/gen/go v0.0.0
	github.com/Glistand/HelpDesk/libs/grpckit v0.0.0
	github.com/Glistand/HelpDesk/libs/otelkit v0.0.0
	github.com/jackc/pgx/v5 v5.7.2
	google.golang.org/grpc v1.83.2
)

replace (
	github.com/Glistand/HelpDesk/api/gen/go => ../../api/gen/go
	github.com/Glistand/HelpDesk/libs/grpckit => ../../libs/grpckit
	github.com/Glistand/HelpDesk/libs/otelkit => ../../libs/otelkit
)
