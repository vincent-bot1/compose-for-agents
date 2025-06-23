import os
from typing import Any, AsyncIterable
from a2a.server.apps import A2AStarletteApplication
from a2a.types import (
    AgentCapabilities,
    AgentCard,
)
from a2a.server.request_handlers import DefaultRequestHandler
from a2a.server.tasks import InMemoryTaskStore
from starlette.applications import Starlette
import yaml

from google.adk.agents.llm_agent import LlmAgent
from google.adk.models.base_llm import BaseLlm
from google.adk.models.lite_llm import LiteLlm
from google.adk.artifacts import InMemoryArtifactService
from google.adk.memory.in_memory_memory_service import InMemoryMemoryService
from google.adk.runners import Runner
from google.adk.sessions import InMemorySessionService
from google.genai import types

from AgentKit.tools.mcp import create_mcp_toolsets

from .config import AgentConfig
from .base_agent import BaseAgent
from ..executor import ADKAgentExecutor

SUPPORTED_CONTENT_TYPES = ["text", "text/plain"]


class Agent(BaseAgent):
    """
    Base class for all agents.
    """

    def __init__(self, config: AgentConfig):
        self._config = config
        self._agent = self.build_agent()
        self._user_id = config.agent_id
        self._runner = Runner(
            app_name=self._config.agent_id,
            agent=self._agent,
            artifact_service=InMemoryArtifactService(),
            session_service=InMemorySessionService(),
            memory_service=InMemoryMemoryService(),
        )

    def app(self, port: int) -> Starlette:
        request_handler = DefaultRequestHandler(
            agent_executor=ADKAgentExecutor(self),
            task_store=InMemoryTaskStore(),
        )

        capabilities = AgentCapabilities(streaming=True)

        agent_card = AgentCard(
            name=self._config.name,
            capabilities=capabilities,
            description=self._config.description or "",
            skills=self._config.skills or [],
            url=f"http://0.0.0.0:{port}/",
            version="1.0.0",
            defaultInputModes=SUPPORTED_CONTENT_TYPES,
            defaultOutputModes=SUPPORTED_CONTENT_TYPES,
        )

        server = A2AStarletteApplication(
            agent_card=agent_card,
            http_handler=request_handler,
        )

        return server.build()

    def get_processing_message(self) -> str:
        return "Processing the reimbursement request..."

    def _build_model(self) -> BaseLlm:
        """Builds the LLM model."""
        provider: str | None
        if isinstance(self._config.model, str):
            provider = "docker"
            name = self._config.model
        else:
            provider = self._config.model.provider
            name = self._config.model.name

        if not provider:
            provider = "docker"

        base_url = None
        api_key: str | None = None
        if provider == "docker":
            api_key = "does_not_matter_but_cannot_be_empty"
            base_url = "http://localhost:12434/engines/v1"
        elif provider == "openai":
            api_key = os.getenv("OPENAI_API_KEY")
        else:
            raise ValueError(f"unknown model provider {provider}")
        return LiteLlm(model="openai/" + name, api_key=api_key, base_url=base_url)

    def build_agent(self) -> LlmAgent:
        """Builds the LLM agent."""
        tools= create_mcp_toolsets(tools_cfg=self._config.tools or [])
        return LlmAgent(
            model=self._build_model(),
            name=self._config.agent_id,
            description=self._config.description or "",
            instruction=self._config.instructions or "",
            tools=tools, # type: ignore
        )

    async def stream(
        self, query: str, session_id: str
    ) -> AsyncIterable[dict[str, Any]]:
        session = await self._runner.session_service.get_session(
            app_name=self._agent.name,
            user_id=self._user_id,
            session_id=session_id,
        )
        content = types.Content(role="user", parts=[types.Part.from_text(text=query)])
        if session is None:
            session = await self._runner.session_service.create_session(
                app_name=self._agent.name,
                user_id=self._user_id,
                state={},
                session_id=session_id,
            )
        async for event in self._runner.run_async(
            user_id=self._user_id, session_id=session.id, new_message=content
        ):
            if event.is_final_response():
                response = ""
                if (
                    event.content
                    and event.content.parts
                    and event.content.parts[0].text
                ):
                    response = "\n".join(
                        [p.text for p in event.content.parts if p.text]
                    )
                elif (
                    event.content
                    and event.content.parts
                    and any([True for p in event.content.parts if p.function_response])
                ):
                    response = next(
                        p.function_response.model_dump() for p in event.content.parts if p.function_response
                    ) # type: ignore
                yield {
                    "is_task_complete": True,
                    "content": response,
                }
            else:
                yield {
                    "is_task_complete": False,
                    "updates": self.get_processing_message(),
                }

    @classmethod
    def from_yaml_filename(cls, filename: str) -> "Agent":
        """
        Create an agent instance from a YAML configuration file.
        This method should be overridden by subclasses if they have specific configuration needs.
        """
        with open(filename, "r") as file:
            config_data = yaml.safe_load(file)
        config = AgentConfig(**config_data)
        return cls(config)

    def __str__(self):
        return f"Agent(name={self._config.name})"
