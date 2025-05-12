# syntax=docker/dockerfile:1

FROM golang:1.24-alpine3.21@sha256:ef18ee7117463ac1055f5a370ed18b8750f01589f13ea0b48642f5792b234044 AS build_gateway
WORKDIR /app
RUN --mount=type=cache,target=/root/.cache/go-build \
    --mount=source=.,target=. \
    go build -o / ./cmd/client_gateway/ ./cmd/agents_gateway/

FROM scratch AS client_gateway
ENTRYPOINT ["/client_gateway"]
COPY --from=build_gateway /client_gateway /

FROM alpine:3.21@sha256:a8560b36e8b8210634f77d9f7f9efd7ffa463e380b75e2e74aff4511df3ef88c AS agents_gateway
RUN apk add --no-cache docker-cli cosign
ENTRYPOINT ["/agents_gateway"]
COPY --from=build_gateway /agents_gateway /
