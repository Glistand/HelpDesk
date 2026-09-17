# syntax=docker/dockerfile:1
FROM golang:1.26-alpine AS build
WORKDIR /src
RUN apk add --no-cache git ca-certificates
ENV GOWORK=off
COPY api/gen/go ./api/gen/go
COPY libs ./libs
COPY services/search-service ./services/search-service
WORKDIR /src/services/search-service
RUN go mod download
RUN CGO_ENABLED=0 go build -o /out/search-service ./cmd/search-service

FROM alpine:3.21
RUN apk add --no-cache ca-certificates
COPY --from=build /out/search-service /usr/local/bin/search-service
EXPOSE 50057
ENTRYPOINT ["search-service"]
