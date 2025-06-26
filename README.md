# Compose agents demo

## Prerequisites

+ Install Docker Desktop `4.42.1`

## Demos

Each of these demos is self-contained and can be run either locally or using a cloud context. They are all configured using two steps.

1. change directory to the root of the demo project
1. create a `.mcp.env` file from the `mcp.env.example` file (if it exists, otherwise the demo doesn't need any secrets) and supply the required MCP tokens
2. run `docker compose up --build`

| Demo | Models | MCPs | project | compose |
| ---- | ---- | ---- | ---- | ---- |
| [Agno](https://github.com/agno-agi/agno) agent that summarizes GitHub issues | deepseek(local), qwen3(local), o3(openai) | github-official, notion, fetch | [./agno](./agno) | [compose.yaml](./agno/compose.yaml) |
| [Vercel AI-SDK](https://github.com/vercel/ai) Chat-UI for mixing MCPs and Model | llama3.2(local), qwen3(local) | wikipedia-mcp, brave, resend(email) | [./vercel](./vercel) | [compose.yaml](https://github.com/slimslenderslacks/scira-mcp-chat/blob/main/compose.yaml) |
| [CrewAI](https://github.com/crewAIInc/crewAI) Marketing Strategy Agent | qwen3(local) | duckduckgo | [./crew-ai](./crew-ai) | [compose.yaml](https://github.com/docker/compose-agents-demo/blob/main/crew-ai/compose.yaml) |
| [ADK](https://github.com/google/adk-python) academic_research agent | gemma3-qat(local) | duckduckgo | [./adk](./adk) | [compose.yaml](./adk/compose.yaml) | 
| [LangGraph](https://github.com/langchain-ai/langgraph) SQL Agent | qwen3(local) | postgres | [./langgraph](./langgraph) | [compose.yaml](./langgraph/compose.yaml) |
| [Embabel](https://github.com/embabel/embabel-agent) Travel Agent | qwen3, Claude3.7, llama3.2, jimclark106/all-minilm:23M-F16 | brave, github-official, wikipedia-mcp, weather, google-maps, airbnb | [./embabel](./embabel) | [compose.yaml](https://github.com/embabel/travel-planner-agent/blob/main/compose.yaml) and [compose.dmr.yaml](https://github.com/embabel/travel-planner-agent/blob/main/compose.dmr.yaml) |
| [Spring AI](https://spring.io/projects/spring-ai) Brace Search | none | brave | [./spring-ai](./spring-ai) | [compose.yaml](./spring-ai/compose.yaml) |

* the embabel demo merges two compose files at runtime. This is because it also supports an Ollama configuration so there is a compose.ollama.yaml file and a compose.dmr.yaml file.

- [ ] hard to demo Agno with only local qwen3 models
- [ ] Setting the Vercel AI-SDK demo without secret support makes the MCPs seem not very "dynamic". They're harder to configure now that we can't run `docker mcp secret ...` to configure them.
- [ ] CrewAI demo does not have a UI right now. When you run `docker compose up` the crew will design a marketing strategy for a hard coded domain (described in the src/main.py file).
