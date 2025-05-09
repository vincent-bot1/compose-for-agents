import asyncio
import os

import nest_asyncio
from agno.agent import Agent
from agno.team import Team
from agno.models.openai import OpenAIChat
from agno.playground import Playground, serve_playground_app
from agno.tools.mcp import MCPTools
from fastapi.middleware.cors import CORSMiddleware

# Allow nested event loops
nest_asyncio.apply()

async def run_server() -> None:
    """Run the GitHub agent server."""
    # Create a client session to connect to the MCP server
    async with MCPTools(
        transport="sse", url=f"http://{os.environ['MCPGATEWAY_ENDPOINT']}/sse"
    ) as mcp_tools:
        gemma_model = OpenAIChat(
                id="ai/gemma3",
                base_url="http://model-runner.docker.internal/engines/llama.cpp/v1",
            )
        gemma_model.role_map = {
            "system": "system",
            "user": "user",
            "assistant": "assistant",
            "tool": "tool",
            "model": "assistant",
        }

        # Create individual specialized agents
        researcher = Agent(
            name="Researcher",
            role="Expert at finding information",
            tools=[mcp_tools],
            model=OpenAIChat("gpt-4o"),
        )

        # Create individual specialized agents
        github = Agent(
            name="Github agent",
            description="A specialized agent for GitHub tasks",
            tools=[mcp_tools],
            model=OpenAIChat("gpt-4o"),
        )

        writer = Agent(
            name="Writer",
            role="Expert at writing clear, engaging content",
            model=gemma_model,
        )

        # Create a team with these agents
        content_team = Team(
            name="Content Team",
            mode="coordinate",
            members=[researcher, writer],
            instructions="You are a team of researchers and writers that work together to create high-quality content.",
            model=OpenAIChat("gpt-4o"),
            markdown=True,
        )

        playground = Playground(agents=[github], teams=[content_team])
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


if __name__ == "__main__":
    asyncio.run(run_server())
