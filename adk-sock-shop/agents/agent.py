"""Agent for reviewing new suppliers before adding them to the Sock Store."""

import os

from google.adk import Agent
from google.adk.agents import SequentialAgent
from google.adk.models.lite_llm import LiteLlm

from .sub_agents.reddit_researcher import reddit_researcher_agent
from .sub_agents.customer_feedback import customer_feedback_agent
from .tools import create_mcp_toolsets

catalog_agent = Agent(
        name="catalog_agent",
        model=LiteLlm(model="openai/gpt-4", api_base="https://api.openai.com/v1", api_key=os.environ.get('OPENAI_API_KEY')),
        instruction = """
        You are reviewing new suppliers for whether they should be added to the store or not.
        If you don't think that supplier will be a good fit, then reject them but if you know their email address, then send them an email to let them know they've been rejected and why.
        If you think that supplier is a good fit, then go ahead and approve them, and add a sku to the catalog using our api by using curl to make a POST request to 
        the endpoint http://localhost:8081/catalogue with content type application/json and a payload
        that matches the following example.
        
        ```
        {
          "name": "Not a sock",
          "description": "A dog not a sock",
          "imageUrl": ["https://tinyurl.com/5n6spnvu", "https://tinyurl.com/mv8ebjnh"],
          "price": 12.99,
          "count": 42,
          "tag": ["animal"]
        }
        ```
        Fill out the values of this payload with data from the supplier.
        """,
        tools = create_mcp_toolsets(tools_cfg=["mcp/resend:send-email", "mcp/curl:curl"])
        )

new_supplier_agent = SequentialAgent(
    name="new_supplier_agent",
    description=
        """
        Supplier Intake Agent
        """,
    sub_agents=[reddit_researcher_agent, customer_feedback_agent, catalog_agent],
)

root_agent = new_supplier_agent
