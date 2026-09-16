# syntax=docker/dockerfile:1
FROM golang:1.26-alpine AS build
WORKDIR /src
RUN apk add --no-cache git ca-certificates
ENV GOWORK=off
COPY api/gen/go ./api/gen/go
COPY libs ./libs
COPY services/assignment-service ./services/assignment-service
WORKDIR /src/services/assignment-service
RUN go mod download
RUN CGO_ENABLED=0 go build -o /out/assignment-service ./cmd/assignment-service

FROM alpine:3.21
RUN apk add --no-cache ca-certificates
COPY --from=build /out/assignment-service /usr/local/bin/assignment-service
EXPOSE 50053
ENTRYPOINT ["assignment-service"]
