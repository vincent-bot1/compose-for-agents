# Compose agents demo

Requirements:

- A recent Docker Desktop with the Docker MCP Toolkit extension
- In the MCP Toolkit, configure and enable GitHub Official, DuckDuckGo and SQLite.

Build `docker compose` from `main`:

```console
$ task build-compose
```

Build a gateway compose provider:

```console
$ task -d gateway build-compose-provider
```

Then you can run:

```console
$ docker compose up --build
$ docker compose down --remove-orphans
```

You can the see the agent UI on http://localhost:3000
