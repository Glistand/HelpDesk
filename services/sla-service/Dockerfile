# syntax=docker/dockerfile:1
FROM golang:1.26-alpine AS build
WORKDIR /src
RUN apk add --no-cache git ca-certificates
ENV GOWORK=off
COPY api/gen/go ./api/gen/go
COPY libs ./libs
COPY services/sla-service ./services/sla-service
WORKDIR /src/services/sla-service
RUN go mod download
RUN CGO_ENABLED=0 go build -o /out/sla-service ./cmd/sla-service

FROM alpine:3.21
RUN apk add --no-cache ca-certificates
COPY --from=build /out/sla-service /usr/local/bin/sla-service
EXPOSE 50055
ENTRYPOINT ["sla-service"]
