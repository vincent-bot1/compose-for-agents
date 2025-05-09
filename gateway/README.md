## How to get it running?

Build the `docker/mcpgateway` image:

```console
task build-gateway-image
```

Build the `docker-mcp` compose provider and install it as a cli plugin.

```console
task build-compose-provider
```

Manually build the latest version of [docker-compose](https://github.com/docker/compose)
and install it as a cli plugin.

Run the whole stack:

```console
docker compose up --build
docker compose down --remove-orphans
```

## Do it with a single command

```console
task
```
