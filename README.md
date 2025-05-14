# Compose agents demo

Requirements:

+ Make sure you have Docker Desktop 4.41 or a recent 4.42 nightly.
+ Install the MCP Toolkit extension (version 1.0.0 on DD 4.41 or version 1.0.1 on 4.42).
+ Configure 3 or 4 MCP Servers in the extension
  + GitHub Official <-- needs a token
  + DuckDuckGo
  + SQLite
  + (Notion <-- needs a token)

Build and install `docker compose` from `main` and build the all the binaries and images, install the plugins:

```console
$ task install
```

Add your OpenAI API key to your environment:

```console
export OPENAI_API_KEY=<KEY>
```

Then you can run:

```console
$ task up
$ task down
```

Cleanup:

```console
$ task uninstall
```

You can then see the agent UI on http://localhost:3000
