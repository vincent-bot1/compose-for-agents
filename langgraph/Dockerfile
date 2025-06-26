FROM alpine:3.22 AS importer
RUN apk add --no-cache pgloader
COPY <<EOF /entrypoint.sh
#!/bin/sh
set -xe
pgloader "\${SQLITE_FILE}" "\${DATABASE_URL}"
EOF
RUN chmod +x /entrypoint.sh
ENTRYPOINT [ "/entrypoint.sh" ]

FROM python:3.13-slim AS agent
ENV PYTHONUNBUFFERED=1
RUN pip install uv

WORKDIR /app
COPY pyproject.toml uv.lock ./
RUN --mount=type=cache,target=/root/.cache/uv \
    UV_COMPILE_BYTECODE=1 uv pip install --system .
COPY agent.py .
RUN python -m compileall -q .
ENTRYPOINT [ "python", "agent.py" ]
