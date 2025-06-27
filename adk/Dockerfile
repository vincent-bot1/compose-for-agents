# Use Python 3.11 slim image as base
FROM python:3.13-slim
ENV PYTHONUNBUFFERED=1

RUN pip install uv

WORKDIR /app
# Install system dependencies
COPY pyproject.toml uv.lock ./
RUN --mount=type=cache,target=/root/.cache/uv \
    UV_COMPILE_BYTECODE=1 uv pip install --system .
# Copy application code
COPY agents/ ./agents/
RUN python -m compileall -q .

# Create non-root user
RUN useradd --create-home --shell /bin/bash app \
    && chown -R app:app /app
USER app

CMD [ "adk", "web", "--host", "0.0.0.0", "--port", "8080", "--log_level", "DEBUG" ]
