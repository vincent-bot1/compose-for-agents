import asyncio
from os import getenv
import os
from textwrap import dedent

import nest_asyncio
from agno.agent import Agent
from agno.models.openai import OpenAIChat
from agno.playground import Playground, serve_playground_app
from agno.storage.agent.sqlite import SqliteAgentStorage
from agno.tools.mcp import MCPTools
from fastapi.middleware.cors import CORSMiddleware

# Allow nested event loops
nest_asyncio.apply()

agent_storage_file: str = "tmp/agents.db"

async def run_server() -> None:
    """Run the GitHub agent server."""
    # Create a client session to connect to the MCP server
    async with MCPTools(
        transport="sse", url=f"http://{os.environ['MCPGATEWAY_HOST']}/sse"
    ) as mcp_tools:
        agent = Agent(
            name="MCP GitHub Agent",
            tools=[mcp_tools],
            instructions=dedent("""\
                You are a GitHub assistant. Help users explore repositories and their activity.

                - Use headings to organize your responses
                - Be concise and focus on relevant information\
            """),
            model=OpenAIChat(id="gpt-4o"),
            storage=SqliteAgentStorage(
                table_name="basic_agent",
                db_file=agent_storage_file,
                auto_upgrade_schema=True,
            ),
            add_history_to_messages=True,
            num_history_responses=3,
            add_datetime_to_instructions=True,
            markdown=True,
        )

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

        gemma = Agent(
            name="Gemma",
            instructions=dedent("""\
                You are a GitHub assistant. Help users explore repositories and their activity.

                - Use headings to organize your responses
                - Be concise and focus on relevant information\
            """),
            model=gemma_model,
            add_history_to_messages=False,
        )

        playground = Playground(agents=[agent, gemma])
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
