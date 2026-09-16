# syntax=docker/dockerfile:1
FROM golang:1.26-alpine AS build
WORKDIR /src
RUN apk add --no-cache git ca-certificates
ENV GOWORK=off
COPY api/gen/go ./api/gen/go
COPY libs ./libs
COPY services/ticket-service ./services/ticket-service
WORKDIR /src/services/ticket-service
RUN go mod download
RUN CGO_ENABLED=0 go build -o /out/ticket-service ./cmd/ticket-service

FROM alpine:3.21
RUN apk add --no-cache ca-certificates
COPY --from=build /out/ticket-service /usr/local/bin/ticket-service
EXPOSE 50052
ENTRYPOINT ["ticket-service"]
