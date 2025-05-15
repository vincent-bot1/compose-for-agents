# Compose agents demo

## Requirements:

+ Install Docker Desktop `4.41` or a recent `4.42 nightly`.
+ Install the MCP Toolkit extension (version `1.0.0` on DD `4.41` or version `1.0.1` on `4.42`).
+ Configure 3 or 4 MCP Servers in the extension
  + `GitHub Official` <-- needs a token
  + `DuckDuckGo`
  + `SQLite`
  + (`Notion` <-- needs a token)

## Build images and docker CLI plugins

Either with `task install` or manually:

On Mac:

```console
docker buildx bake --file docker-bake.hcl darwin
ln -sf ./bin/docker-compose ./bin/docker-mcpgateway ~/.docker/cli-plugins/
```

On Windows:

```console
docker buildx bake --file docker-bake.hcl windows
mkdir -p "%USERPROFILE%/.docker/cli-plugins"
ln -sf ./bin/docker-compose.exe "%USERPROFILE%/.docker/cli-plugins/"
ln -sf ./bin/docker-mcpgateway.exe "%USERPROFILE%/.docker/cli-plugins/"
```

## Prepare for the run

Add your OpenAI API key to your environment:

```console
export OPENAI_API_KEY=<KEY>
```

## And run!

Start the compose file:

```console
docker compose up
docker compose down --remove-orphans
```

**You can then see the agent UI on http://localhost:3000**

As an alternative, you can also use`task`

```console
task up
task down
```

Cleanup:

```console
task uninstall
```