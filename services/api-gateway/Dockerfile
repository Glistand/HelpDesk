# syntax=docker/dockerfile:1
FROM golang:1.26-alpine AS build
WORKDIR /src
RUN apk add --no-cache git ca-certificates
ENV GOWORK=off
COPY api/gen/go ./api/gen/go
COPY libs ./libs
COPY services/api-gateway ./services/api-gateway
WORKDIR /src/services/api-gateway
RUN go mod download
RUN CGO_ENABLED=0 go build -o /out/api-gateway ./cmd/api-gateway

FROM alpine:3.21
RUN apk add --no-cache ca-certificates
COPY --from=build /out/api-gateway /usr/local/bin/api-gateway
EXPOSE 8080
ENTRYPOINT ["api-gateway"]
