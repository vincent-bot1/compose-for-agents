# Compose agents demo

## Prerequisites

+ Install Docker Desktop `4.41` or a recent `4.42 nightly`
+ If using 4.41, install the MCP Toolkit extension (version `1.0.0`)
+ If using a 4.42 nightly, use the MCP Toolkit in the left nav

## OpenAI API Key

1. Generate a key by navigating to https://platform.openai.com/api-keys
1. Create a `.env` file and add your OpenAI API key to it:
```console
OPENAI_API_KEY=<KEY>
```

## Configure MCP Servers

+ Configure the following MCP Servers in the Docker Desktop extension (Desktop 4.41) or MCP Toolkit left nav (Desktop 4.42 nightly):
  + Notion
  + GitHub Official

### Notion

1. Create a new Notion account using a non-company email address
1. Create a new integration by navigating to https://www.notion.so/my-integrations
1. Follow the creation flow enabling write access
1. Add the Notion MCP Server in Docker Desktop
![Notion MCP extension](./img/notion-mcp-server.png)
1. Copy the integration token into Desktop's MCP server configuration
![Notion token](./img/notion-token.png)
![Notion MCP config](./img/notion-mcp-config.png)
1. Create a page named "Updates" in your workspace
1. Give your integration access to the page by clicking on the ... menu on the top right of the updates page, clicking "Connections" and selecting it
![Notion page perms](./img/notion-page-perms.png)

### GitHub Official

1. Create a fine grained personal access token: https://github.com/settings/personal-access-tokens
1. Give it read access to public repos
![GitHub token perms](./img/github-perms.png)
1. Add the "GitHub Official" MCP server
![GitHub MCP server](./img/github-mcp-server.png)
1. Add your token to it

## Build Docker CLI Plugin

Either with `task install` or manually:

On Mac:

```console
docker buildx bake --file docker-bake.hcl darwin
ln -sf $(pwd)/bin/docker-compose ~/.docker/cli-plugins/
```

On Windows:

```console
docker buildx bake --file docker-bake.hcl windows
Copy-Item -Path ./bin/docker-compose.exe -Destination "$env:USERPROFILE\.docker\cli-plugins"
```

## And Run!

Start the application:

```console
docker compose up --build
docker compose down --remove-orphans
```

**You can then see the agent UI on http://localhost:3000**

As an alternative, you can also use `task`

```console
task up
task down
```

Cleanup:

```console
task down
```
