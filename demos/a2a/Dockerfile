##### Base image
FROM python:3.13-slim AS base

ENV PYTHONUNBUFFERED=1 \
    UV_COMPILE_BYTECODE=1 \
    UV_NO_CACHE=1

RUN pip install uv

WORKDIR /app
COPY pyproject.toml uv.lock ./
RUN uv pip install --system .
COPY main.py src ./
RUN python -m compileall -q .

RUN useradd --create-home --shell /bin/bash app \
    && chown -R app:app /app
USER app

ENV AGENT_CONFIG=/app/agent.yaml

##### Critic Agent
FROM base AS critic-agent

COPY --chown=app:app agents/critic.yaml /app/agent.yaml
EXPOSE 8001
CMD ["python", "main.py", "--host", "0.0.0.0", "--port", "8001"]

##### Auditor Agent
FROM base AS auditor-agent

COPY --chown=app:app agents/auditor.yaml /app/agent.yaml
EXPOSE 8002
CMD ["adk", "web", ".", "--host", "0.0.0.0", "--port", "8002"]
#CMD ["python", "main.py", "--host", "0.0.0.0", "--port", "8002"]
