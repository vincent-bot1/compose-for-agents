import asyncio
import os
import sys
from typing import Optional

import nest_asyncio
import yaml

from agno import agent, team
from agno.models.openai import OpenAIChat
from agno.playground import Playground, serve_playground_app
from agno.tools.mcp import MCPTools, Toolkit
from agno.tools.reasoning import ReasoningTools
from fastapi.middleware.cors import CORSMiddleware

# Allow nested event loops
nest_asyncio.apply()

DOCKER_MODEL_PROVIDER = "docker"

class Agent(agent.Agent):
    @property
    def is_streamable(self) -> bool:
        if self.stream is not None:
            return self.stream
        return super().is_streamable


class Team(team.Team):
    @property
    def is_streamable(self) -> bool:
        stream = getattr(self, "stream")
        if stream is not None:
            return stream
        return super().is_streamable


def should_stream(model_provider: str, tools: list[Toolkit]) -> Optional[bool]:
    """Returns whether a model with the given provider and tools can stream"""
    if model_provider == DOCKER_MODEL_PROVIDER and len(tools) > 0:
        # DMR doesn't yet support tools with streaming
        return True
    # Let the model/options decide
    return None


def create_model(model_name: str, provider: str, temperature: float) -> OpenAIChat:
    """Create a model instance based on the model name and provider."""
    print(f"creating model {model_name} with provider {provider} and temperature {temperature}")
    if provider == DOCKER_MODEL_PROVIDER:
        base_url = os.getenv("LLM_URL")
        if base_url is None:
            base_url = "http://model-runner.docker.internal/engines/llama.cpp/v1"
        model = OpenAIChat(id=model_name, base_url=base_url, temperature=temperature)
        model.role_map = {
            "system": "system",
            "user": "user",
            "assistant": "assistant",
            "tool": "tool",
            "model": "assistant",
        }
        return model

    if provider == "openai":
        if os.environ.get("OPENAI_API_KEY") is None:
            raise ValueError(
                "OPENAI_API_KEY environment variable not set for OpenAI model"
            )
        return OpenAIChat(model_name, temperature=temperature)

    raise ValueError(f"Unknown agent model provider: {provider}")


async def run_server(config) -> None:
    """Run the playground server."""
    # Create a client session to connect to the MCP server
    agents = []
    agents_by_id = {}
    teams = []
    teams_by_id = {}

    for agent_id, agent_data in config.get("agents", {}).items():
        model_name = agent_data.get("model")
        if not model_name:
            raise ValueError(f"Model name not specified for agent {agent_id}")
        temperature = agent_data.get("temperature", None)
        provider = agent_data.get("model_provider", "docker")
        model = create_model(model_name, provider, temperature)
        markdown = agent_data.get("markdown", False)
        tools: list[Toolkit] = [
#            ReasoningTools(think=True, analyze=True)
        ]
        tools_list = agent_data.get("tools", [])
        if len(tools_list) > 0:
            tool_names = [name.split(":", 1)[1] for name in tools_list]

            # Always use socat, but the endpoint can be different (mock vs real gateway)
            endpoint = os.environ['MCPGATEWAY_ENDPOINT']
            print(f"DEBUG: Connecting to MCP gateway at {endpoint}")

            # Parse endpoint to extract host and port
            import socket
            from urllib.parse import urlparse

            try:
                # Handle both URL format (http://host:port/path) and host:port format
                if endpoint.startswith('http://') or endpoint.startswith('https://'):
                    parsed = urlparse(endpoint)
                    host = parsed.hostname
                    port = parsed.port
                    tcp_endpoint = f"{host}:{port}"
                else:
                    # Legacy host:port format
                    host, port = endpoint.split(':')
                    port = int(port)
                    tcp_endpoint = endpoint

                # Test TCP connection first
                sock = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
                sock.settimeout(5)
                sock.connect((host, port))
                sock.close()
                print(f"DEBUG: TCP connection to {host}:{port} successful")
            except Exception as e:
                print(f"ERROR: TCP connection to {endpoint} failed: {e}")
                raise

            t = MCPTools(
                command=f"socat STDIO TCP:{tcp_endpoint}",
                include_tools=tool_names,
            )
            mcp_tools = await t.__aenter__()
            tools = [mcp_tools]
        agent = Agent(
            name=agent_data["name"],
            role=agent_data.get("role", ""),
            description=agent_data.get("description"),
            instructions=agent_data.get("instructions"),
            tools=tools,  # type: ignore,
            model=model,
            show_tool_calls=True,
            stream=should_stream(provider, tools),
            add_datetime_to_instructions=True,
            markdown=markdown,
            debug_mode=True,
        )
        agents_by_id[agent_id] = agent
        # Append only agents that we want to chat with
        if agent_data.get("chat", True):
            agents.append(agent)

    for team_id, team_data in config.get("teams", {}).items():
        model_name = team_data.get("model")
        if not model_name:
            raise ValueError(f"Model name not specified for team {team_id}")
        provider = team_data.get("model_provider", "docker")
        temperature = team_data.get("temperature", None)
        model = create_model(model_name, provider, temperature)
        team_agents: list[Agent | Team] = []
        for agent_id in team_data.get("members", []):
            try:
                agent = agents_by_id[agent_id]
            except KeyError:
                raise ValueError(f"Agent {agent_id} not found in agents")
            team_agents.append(agent)
        markdown = team_data.get("markdown", False)
        team_tools: list[Toolkit] = [
#            ReasoningTools(think=True, analyze=True)
        ]
        tools_list = team_data.get("tools", [])
        if len(tools_list) > 0:
            tool_names = [name.split(":", 1)[1] for name in tools_list]

            # Always use socat, but the endpoint can be different (mock vs real gateway)
            endpoint = os.environ['MCPGATEWAY_ENDPOINT']
            print(f"DEBUG: Team connecting to MCP gateway at {endpoint}")

            # Parse endpoint to extract host and port
            from urllib.parse import urlparse

            # Handle both URL format (http://host:port/path) and host:port format
            if endpoint.startswith('http://') or endpoint.startswith('https://'):
                parsed = urlparse(endpoint)
                host = parsed.hostname
                port = parsed.port
                tcp_endpoint = f"{host}:{port}"
            else:
                # Legacy host:port format
                tcp_endpoint = endpoint

            t = MCPTools(
                command=f"socat STDIO TCP:{tcp_endpoint}",
                include_tools=tool_names,
            )
            mcp_tools = await t.__aenter__()
            team_tools = [mcp_tools]
        team = Team(
            name=team_data.get("name", ""),
            mode=team_data.get("mode", "coordinate"),
            members=team_agents, # type: ignore
            description=team_data.get("description"),
            instructions=team_data.get("instructions"),
            tools=team_tools,  # type: ignore,
            model=model,
            # show_members_responses=True,
            # show_tool_calls=True,
            markdown=markdown,
            add_datetime_to_instructions=True,
            debug_mode=True,
        )
        team.stream = should_stream(provider, team_tools)
        teams_by_id[team_id] = team
        if team_data.get("chat", True):
            teams.append(team)

    playground = Playground(agents=agents, teams=teams) # type: ignore

    app = playground.get_app()
    app.add_middleware(
        CORSMiddleware,
        allow_origins=["*"],
        allow_credentials=True,
        allow_methods=["*"],
        allow_headers=["*"],
    )

    # Serve the app while keeping the MCPTools context manager alive
    serve_playground_app(host="0.0.0.0", app=app)


def main():
    config_filename = sys.argv[1] if len(sys.argv) > 1 else "/agents.yaml"
    with open(config_filename, "r") as f:
        config = yaml.safe_load(f)

    asyncio.run(run_server(config))


if __name__ == "__main__":
    main()
