"""Reddit Researcher Agent"""

import os

from google.adk import Agent
from google.adk.models.lite_llm import LiteLlm

from . import prompt
from ...tools import create_mcp_toolsets

tools = create_mcp_toolsets(tools_cfg=["mcp/brave:brave_web_search"])

reddit_researcher_agent = Agent(
    # Using GPT-4 from OpenAI cloud
    model=LiteLlm(model="openai/gpt-4", api_base="https://api.openai.com/v1", api_key=os.environ.get('OPENAI_API_KEY')),
    name="reddit_researcher_agent",
    instruction=prompt.PROMPT,
    tools=tools,  # type: ignore
)
