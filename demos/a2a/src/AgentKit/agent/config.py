from typing import Optional, Union


from a2a.types import AgentSkill
from pydantic import BaseModel


class ModelSpec(BaseModel):
    """
    Specification for a model used by an agent.
    If the provider is not specified, it defaults to 'docker'.
    """

    name: str
    provider: Optional[str] = None


class AgentConfig(BaseModel):
    """
    Configuration for an agent.
    """

    name: str
    id: Optional[str] = None
    description: Optional[str] = None
    instructions: Optional[str] = None
    model: Union[str, ModelSpec]  # Model can be a string or a ModelSpec object
    skills: Optional[list[AgentSkill]]
    tools: Optional[list[str]] = None

    @property
    def agent_id(self) -> str:
        if self.id:
            return self.id
        return self.name.replace(" ", "_").lower()
