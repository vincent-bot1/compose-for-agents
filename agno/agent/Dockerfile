FROM python:3.13-slim
ENV PYTHONUNBUFFERED=1
RUN pip install uv

WORKDIR /app
COPY pyproject.toml uv.lock ./
RUN --mount=type=cache,target=/root/.cache/uv \
    UV_COMPILE_BYTECODE=1 uv pip install --system .
COPY playground.py .
RUN python -m compileall -q .
COPY <<EOF entrypoint.sh
#!/bin/sh
if test -n "\${OPENAI_API_KEY}"; then
    echo "OPENAI_API_KEY is set, using OpenAI with \${OPENAI_MODEL_NAME}"
    export MODEL_PROVIDER=openai
    export MODEL_NAME=\${OPENAI_MODEL_NAME}
else
    echo "OPENAI_API_KEY is not set, using Docker Model Runner with \${MODEL_RUNNER_MODEL}"
    export MODEL_PROVIDER=docker
    export MODEL_NAME=\${MODEL_RUNNER_MODEL}
    export OPENAI_API_KEY=cannot_be_empty
fi
exec python playground.py
EOF
RUN chmod +x entrypoint.sh
ENTRYPOINT ["./entrypoint.sh"]
