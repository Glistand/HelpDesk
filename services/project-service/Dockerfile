# syntax=docker/dockerfile:1
FROM golang:1.26-alpine AS build
WORKDIR /src
RUN apk add --no-cache git ca-certificates
ENV GOWORK=off
COPY api/gen/go ./api/gen/go
COPY libs ./libs
COPY services/project-service ./services/project-service
WORKDIR /src/services/project-service
RUN go mod download
RUN CGO_ENABLED=0 go build -o /out/project-service ./cmd/project-service

FROM alpine:3.21
RUN apk add --no-cache ca-certificates
COPY --from=build /out/project-service /usr/local/bin/project-service
EXPOSE 50059
ENTRYPOINT ["project-service"]
