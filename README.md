# Compose agents demo

Requirements:

Build and install `docker compose` from `main`:

```console
$ task build-compose
```

Build the all the binaries and images, install the plugins:

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
