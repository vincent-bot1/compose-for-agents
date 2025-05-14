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
    def __init__(self, *args, **kwargs):
        super().__init__(*args, **kwargs)
        self.stream = None

    @property
    def is_streamable(self) -> bool:
        if self.stream is not None:
            return self.stream
        return super().is_streamable


def should_stream(model_provider: str, tools: list[Toolkit]) -> Optional[bool]:
    """Returns whether a model with the given provider and tools can stream"""
    if model_provider == DOCKER_MODEL_PROVIDER and len(tools) > 0:
        # DMR doesn't yet support tools with streaming
        return False
    # Let the model/options decide
    return None


def create_model(model_name: str, provider: str) -> OpenAIChat:
    """Create a model instance based on the model name and provider."""
    if provider == DOCKER_MODEL_PROVIDER:
        base_url = os.getenv("AI_RUNNER_URL")
        if base_url is None:
            base_url = "http://model-runner.docker.internal/engines/llama.cpp/v1"
        model = OpenAIChat(id="ai/" + model_name, base_url=base_url)
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
        return OpenAIChat(model_name)

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
        provider = agent_data.get("model_provider", "docker")
        model = create_model(model_name, provider)
        markdown = agent_data.get("markdown", False)
        tools: list[Toolkit] = []
        tools_list = agent_data.get("tools", [])
        if len(tools_list) > 0:
            tool_names = [name.split(":", 1)[1] for name in tools_list]
            t = MCPTools(
                command=f"socat STDIO TCP:{os.environ['MCPGATEWAY_ENDPOINT']}",
                include_tools=tool_names,
            )
            mcp_tools = await t.__aenter__()
            tools = [mcp_tools]
        agent = Agent(
            name=agent_data["name"],
            role=agent_data.get("role", ""),
            description=agent_data.get("description", ""),
            tools=tools,  # type: ignore,
            model=model,
            stream=should_stream(provider, tools),
            markdown=markdown,
        )
        agents_by_id[agent_id] = agent
        # Append only agents that we want to chat with
        if agent_data.get("chat", True):
            agents.append(agent)

    for team_id, team_data in config.get("teams", {}).items():
        model_name = agent_data.get("model")
        if not model_name:
            raise ValueError(f"Model name not specified for team {team_id}")
        provider = agent_data.get("model_provider", "docker")
        model = create_model(model_name, provider)
        team_agents: list[Agent | Team] = []
        for agent_id in team_data.get("members", []):
            try:
                agent = agents_by_id[agent_id]
            except KeyError:
                raise ValueError(f"Agent {agent_id} not found in agents")
            team_agents.append(agent)
        markdown = agent_data.get("markdown", False)
        team_tools: list[Toolkit] = []
        tools_list = agent_data.get("tools", [])
        if len(tools_list) > 0:
            tool_names = [name.split(":", 1)[1] for name in tools_list]
            t = MCPTools(
                command=f"socat STDIO TCP:{os.environ['MCPGATEWAY_ENDPOINT']}",
                include_tools=tool_names,
            )
            mcp_tools = await t.__aenter__()
            team_tools = [mcp_tools]
        team = Team(
            name=team_data.get("name", ""),
            mode=team_data.get("mode", "coordinate"),
            members=team_agents,
            instructions=team_data.get("instructions", ""),
            tools=team_tools,  # type: ignore,
            model=model,
            markdown=markdown,
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
