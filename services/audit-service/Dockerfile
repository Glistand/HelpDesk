# syntax=docker/dockerfile:1
FROM golang:1.26-alpine AS build
WORKDIR /src
RUN apk add --no-cache git ca-certificates
ENV GOWORK=off
COPY api/gen/go ./api/gen/go
COPY libs ./libs
COPY services/audit-service ./services/audit-service
WORKDIR /src/services/audit-service
RUN go mod download
RUN CGO_ENABLED=0 go build -o /out/audit-service ./cmd/audit-service

FROM alpine:3.21
RUN apk add --no-cache ca-certificates
COPY --from=build /out/audit-service /usr/local/bin/audit-service
EXPOSE 50054
ENTRYPOINT ["audit-service"]
