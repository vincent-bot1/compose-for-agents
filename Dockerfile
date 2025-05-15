# syntax=docker/dockerfile:1

# Build the client image (not used in the main demo)
FROM golang:1.24-alpine3.21@sha256:ef18ee7117463ac1055f5a370ed18b8750f01589f13ea0b48642f5792b234044 AS build_client
WORKDIR /app
RUN --mount=type=cache,target=/root/.cache/go-build,id=client \
    --mount=source=.,target=. \
    go build -o / ./cmd/client/

FROM scratch AS client
ENTRYPOINT ["/client"]
COPY --from=build_client /client /

# Build the agents_gateway image
FROM golang:1.24-alpine3.21@sha256:ef18ee7117463ac1055f5a370ed18b8750f01589f13ea0b48642f5792b234044 AS build_agents_gateway
WORKDIR /app
RUN --mount=type=cache,target=/root/.cache/go-build,id=agents_gateway \
    --mount=source=.,target=. \
    go build -o / ./cmd/agents_gateway/

FROM alpine:3.21@sha256:a8560b36e8b8210634f77d9f7f9efd7ffa463e380b75e2e74aff4511df3ef88c AS agents_gateway
RUN apk add --no-cache docker-cli
ENTRYPOINT ["/agents_gateway"]
COPY --from=build_agents_gateway /agents_gateway /

# Build the docker-mcpgateway compose provider (darwin)
FROM --platform=$BUILDPLATFORM golang:1.24-alpine3.21@sha256:ef18ee7117463ac1055f5a370ed18b8750f01589f13ea0b48642f5792b234044 AS build_docker-mcpgateway-darwin
WORKDIR /app
RUN --mount=type=cache,target=/root/.cache/go-build,id=docker-mcpgateway \
    --mount=source=.,target=. \
    CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -o /out/docker-mcpgateway ./cmd/docker-mcpgateway/

FROM scratch AS docker-mcpgateway-darwin
COPY --from=build_docker-mcpgateway-darwin /out/docker-mcpgateway /

# Build the docker-mcpgateway compose provider (windows)
FROM --platform=$BUILDPLATFORM golang:1.24-alpine3.21@sha256:ef18ee7117463ac1055f5a370ed18b8750f01589f13ea0b48642f5792b234044 AS build_docker-mcpgateway-windows
WORKDIR /app
RUN --mount=type=cache,target=/root/.cache/go-build,id=docker-mcpgateway \
    --mount=source=.,target=. \
    CGO_ENABLED=0 GOOGS=windows GOARCH=amd64 go build -o /out/docker-mcpgateway.exe ./cmd/docker-mcpgateway/

FROM scratch AS docker-mcpgateway-windows
COPY --from=build_docker-mcpgateway-windows /out/docker-mcpgateway.exe /