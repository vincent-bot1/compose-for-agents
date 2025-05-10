## How to get it running?

Build and install the required dependencies

```console
task build_and_install
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
