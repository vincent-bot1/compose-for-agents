# Copyright 2025 Google LLC
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

"""LLM Auditor for verifying & refining LLM-generated answers using the web."""

from google.adk.agents import SequentialAgent

from .sub_agents.reddit_researcher import reddit_researcher_agent
from .sub_agents.user_feedback import user_feedback_agent

new_supplier_agent = SequentialAgent(
    name="new_supplier_agent",
    description=(
        """
        You are reviewing new suppliers for whether they should be added to the store or not.
        If you are unsure, ask for more information.
        If you don't think the supplier will be a good fit for the store, then reject them.
        However, if you think that the supplier will be a good fit, then go ahead and approve them.
        """
    ),
    sub_agents=[reddit_researcher_agent, user_feedback_agent],
)

root_agent = new_supplier_agent
