# Compose agents demo

Requirements:

- A recent compose version (build from main)
- A recent Docker Desktop with the Docker MCP Toolkit extension
- In the MCP Toolkit, configure and enable GitHub Official, DuckDuckGo and SQLite.

Then you can run:

```console
$ task -d gateway build-compose-provider
$ docker compose up --build
```

You can the see the agent UI on http://localhost:3000
