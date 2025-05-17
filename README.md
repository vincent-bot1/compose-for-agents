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
ln -sf $(pwd)/bin/docker-compose $(pwd)/bin/docker-mcpgateway ~/.docker/cli-plugins/
```

On Windows:

```console
docker buildx bake --file docker-bake.hcl windows
Copy-Item -Path ./bin/docker-compose.exe -Destination "$env:USERPROFILE\.docker\cli-plugins"
Copy-Item -Path ./bin/docker-mcpgateway.exe -Destination "$env:USERPROFILE\.docker\cli-plugins"
```

## Prepare Notion

The Notion Page Creator agent will create the pages under another page
titled "Updates". Create this page in the Notion workspace and give access
to the integration (via the ... menu at the top right).

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
