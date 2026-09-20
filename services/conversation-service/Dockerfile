# syntax=docker/dockerfile:1
FROM golang:1.26-alpine AS build
WORKDIR /src
RUN apk add --no-cache git ca-certificates
ENV GOWORK=off
COPY api/gen/go ./api/gen/go
COPY libs ./libs
COPY services/conversation-service ./services/conversation-service
WORKDIR /src/services/conversation-service
RUN go mod download
RUN CGO_ENABLED=0 go build -o /out/conversation-service ./cmd/conversation-service

FROM alpine:3.21
RUN apk add --no-cache ca-certificates
COPY --from=build /out/conversation-service /usr/local/bin/conversation-service
EXPOSE 50058
ENTRYPOINT ["conversation-service"]
