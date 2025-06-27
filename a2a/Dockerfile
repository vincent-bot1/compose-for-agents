##### Base image
FROM python:3.13-slim AS base

ENV PYTHONUNBUFFERED=1
RUN pip install uv

WORKDIR /app
COPY pyproject.toml uv.lock ./
RUN --mount=type=cache,target=/root/.cache/uv \
    UV_COMPILE_BYTECODE=1 uv pip install --system .
COPY main.py src ./
RUN python -m compileall -q .

ENV AGENT_CONFIG=/app/agent.yaml

COPY <<EOF ./entrypoint.sh
#!/bin/bash
set -e

export LLM_AGENT_API_URL=\${DOCKER_MODEL_RUNNER_URL}
export LLM_AGENT_MODEL_NAME=\${DOCKER_MODEL_RUNNER_MODEL}
exec \$@
EOF
RUN chmod +x ./entrypoint.sh

RUN useradd --create-home --shell /bin/bash app \
    && chown -R app:app /app
USER app

##### Critic Agent
FROM base AS critic-agent

COPY --chown=app:app agents/critic.yaml /app/agent.yaml
EXPOSE 8001
CMD ["./entrypoint.sh", "python", "main.py", "--host", "0.0.0.0", "--port", "8001"]

##### Auditor Agent
FROM base AS auditor-agent

COPY --chown=app:app agents/auditor.yaml /app/agent.yaml
EXPOSE 8002
CMD ["./entrypoint.sh", "adk", "web", ".", "--host", "0.0.0.0", "--port", "8080"]
# Use this to expose the agent as a web service instead of a UI
#CMD ["python", "main.py", "--host", "0.0.0.0", "--port", "8002"]
