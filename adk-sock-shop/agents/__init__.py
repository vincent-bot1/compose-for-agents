"""Supplier Intake Agent."""

import logging
import os

import litellm

from . import agent

# Set the base URL for the OpenAI API to the Docker Model Runner URL
os.environ.setdefault("OPENAI_BASE_URL", os.getenv("MODEL_RUNNER_URL", ""))
os.environ["OTEL_SDK_DISABLED"] = "true"

# Enable logging with reduced verbosity
logging.basicConfig(
    level=logging.INFO,  # Less verbose than DEBUG
    format="%(asctime)s - %(name)s - %(levelname)s - %(message)s",
    force=True,  # override ADK defaults
)
logging.getLogger("opentelemetry").setLevel(logging.ERROR)
logging.getLogger("google.adk").setLevel(logging.INFO)
logging.getLogger("LiteLLM").setLevel(logging.WARNING)  # Much less verbose
logging.getLogger("litellm").setLevel(logging.WARNING)  # Also reduce this
logging.getLogger("httpx").setLevel(logging.WARNING)  # Reduce HTTP logs
logging.getLogger("httpcore").setLevel(logging.WARNING)  # Reduce HTTP core logs
litellm.set_verbose = False  # Disable raw HTTP logs # type: ignore

__all__ = ["agent"]
