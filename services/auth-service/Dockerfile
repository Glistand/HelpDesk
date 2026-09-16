# syntax=docker/dockerfile:1
FROM golang:1.26-alpine AS build
WORKDIR /src
RUN apk add --no-cache git ca-certificates
ENV GOWORK=off
COPY api/gen/go ./api/gen/go
COPY libs ./libs
COPY services/auth-service ./services/auth-service
WORKDIR /src/services/auth-service
RUN go mod download
RUN CGO_ENABLED=0 go build -o /out/auth-service ./cmd/auth-service

FROM alpine:3.21
RUN apk add --no-cache ca-certificates
COPY --from=build /out/auth-service /usr/local/bin/auth-service
EXPOSE 50051
ENTRYPOINT ["auth-service"]
