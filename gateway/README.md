## How to get it running?

**Requirements:**

+ Make sure you have Docker Desktop 4.41 or a recent 4.42 nightly.
+ Install the MCP Toolkit extension (version 1.0.0 on DD 4.41 or version 1.0.1 on 4.42).
+ Configure 3 or 4 MCP Servers in the extension
  + GitHub Official <-- needs a token
  + DuckDuckGo
  + SQLite
  + (Notion <-- needs a token)

**Build and install the required dependencies**

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
