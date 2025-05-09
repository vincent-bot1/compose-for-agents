# Compose agents demo

Requirements:

Build and install `docker compose` from `main`:

```console
$ task build-compose
```

Build the MCP gateway provider and container:

```console
$ task -d gateway build
```

Make sure you have a GitHub token in your env:

```console
export GITHUB_TOKEN=<TOKEN>
```

Then you can run:

```console
$ docker compose up --build
$ docker compose down --remove-orphans
```

You can then see the agent UI on http://localhost:3000
