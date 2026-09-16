# syntax=docker/dockerfile:1
FROM golang:1.26-alpine AS build
WORKDIR /src
RUN apk add --no-cache git ca-certificates
ENV GOWORK=off
COPY libs ./libs
COPY services/notification-service ./services/notification-service
WORKDIR /src/services/notification-service
RUN go mod download
RUN CGO_ENABLED=0 go build -o /out/notification-service ./cmd/notification-service

FROM alpine:3.21
RUN apk add --no-cache ca-certificates
COPY --from=build /out/notification-service /usr/local/bin/notification-service
ENTRYPOINT ["notification-service"]
