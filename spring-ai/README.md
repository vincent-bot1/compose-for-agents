# 🧠 Spring AI + DuckDuckGo with Model Context Protocol (MCP)

This project demonstrates a **zero-config Spring Boot application** using [Spring AI] and the **Model Context Protocol (MCP)** to answer natural language questions by performing real-time web search via [DuckDuckGo] — all orchestrated with [Docker Compose].

> [!Tip]
> ✨ No configuration needed — run it with a single command.

<p align="center">
  <img src="demo.gif"
       alt="Spring AI DuckDuckGo Search Demo"
       width="500"
       style="border: 1px solid #ccc; border-radius: 8px;" />
</p>


# 🚀 Getting Started

### Requirements

- ✅ [Docker Desktop] **v4.43.0** or later

### Run the project

```sh
docker compose up
```

No setup, API keys, or additional configuration required.

# ❓ What Can It Do?

Ask natural language questions and let Spring AI + Brave Search provide intelligent, real-time answers:

- “Does Spring AI support the Model Context Protocol?”
- “What is the Brave Search API?”
- “Give me examples of Spring Boot AI integrations.”

The application uses:
- 	An MCP-compatible gateway to route queries to DuckDuckGo Search
- Spring AI’s LLM client to embed results into answers
- Auto-configuration via Spring Boot to bind everything

You can **customize the question** asked to the agent — just edit the question in `compose.yaml`.

# 🧱 Project Structure

| **File/Folder**          | **Purpose**                                      |
| ------------------------ | ------------------------------------------------ |
| `compose.yml`            | launches the Brave MCP gateway and Spring AI app |
| `Dockerfile`             | Builds the Spring Boot container                 |
| `application.properties` | Sets the MCP gateway URL used by Spring AI       |
| `Application.java`       | Configures the ChatClient with MCP and runs it   |
| `mvnw`, `pom.xml`        | Maven wrapper and build definition               |

# 🔧 Architecture Overview

```mermaid

flowchart TD
    A[($QUESTION)] --> B[Spring Boot App]
    B --> C[Spring AI ChatClient]
    C -->|uses| D[MCP Tool Callback]
    D -->|queries| E[Docker MCP Gateway]
    E -->|calls| F[DuckDuckGo Search API]
    F --> E --> D --> C
    C -->|LLM| H[(Docker Model Runner)]
    H --> C
    C --> G[Final Answer]

```

- The application loads a a question via the `QUESTION` environment variable.
- MCP is used as a tool in the LLM pipeline.
- The response is enriched with real-time DuckDuckGo Search results.

# 📎 Credits

- [Spring AI]
- [DuckDuckGo]
- [Docker Compose]

[DuckDuckGo]: https://duckduckgo.com
[Spring AI]: https://github.com/spring-projects/spring-ai
[Docker Compose]: https://docs.docker.com/compose/
[Docker Desktop]: https://www.docker.com/products/docker-desktop/
