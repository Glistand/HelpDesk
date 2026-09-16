# syntax=docker/dockerfile:1
FROM golang:1.26-alpine AS build
WORKDIR /src
RUN apk add --no-cache git ca-certificates
ENV GOWORK=off
COPY api/gen/go ./api/gen/go
COPY libs ./libs
COPY services/escalation-service ./services/escalation-service
WORKDIR /src/services/escalation-service
RUN go mod download
RUN CGO_ENABLED=0 go build -o /out/escalation-service ./cmd/escalation-service

FROM alpine:3.21
RUN apk add --no-cache ca-certificates
COPY --from=build /out/escalation-service /usr/local/bin/escalation-service
EXPOSE 50056
ENTRYPOINT ["escalation-service"]
