# syntax=docker/dockerfile:1

FROM golang:1.24-alpine3.21@sha256:ef18ee7117463ac1055f5a370ed18b8750f01589f13ea0b48642f5792b234044 AS build
WORKDIR /app
RUN --mount=type=cache,target=/root/.cache/go-build \
    --mount=source=.,target=. \
    go build -o / ./cmd/agent/ ./cmd/agents_gateway/

FROM scratch AS agent
ENTRYPOINT ["/agent"]
COPY --from=build /agent /

FROM scratch AS agents_gateway
ENTRYPOINT ["/agents_gateway"]
COPY --from=build /agents_gateway /
