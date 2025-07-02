FROM python:3.12-slim

RUN apt-get update && apt-get install -y build-essential && \
    rm -rf /var/lib/apt/lists/*
RUN pip install poetry
WORKDIR /app
COPY pyproject.toml poetry.lock ./
RUN poetry install --no-interaction --no-root
COPY . .
RUN poetry install
COPY <<EOF /entrypoint.sh
#!/bin/sh
if test -n "\${OPENAI_API_KEY}"; then
    echo "OPENAI_API_KEY is set, using OpenAI with \${MODEL_NAME}"
    export OPENAI_MODEL_NAME=\${MODEL_NAME}
else
    echo "OPENAI_API_KEY is not set, using Docker Model Runner with \${MODEL_RUNNER_MODEL}"
    export OPENAI_BASE_URL=\${MODEL_RUNNER_URL}
    export OPENAI_MODEL_NAME=\${MODEL_RUNNER_MODEL}
    export OPENAI_API_KEY=cannot_be_empty
fi
exec poetry run marketing_posts
EOF
RUN chmod +x /entrypoint.sh
CMD ["/entrypoint.sh"]
