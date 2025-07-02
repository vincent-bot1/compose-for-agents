build-image:
	docker buildx build --builder hydrobuild --platform linux/amd64,linux/arm64 --tag jimclark106/supplier_agent:latest --push --file Dockerfile .

gateway-secrets:
	docker mcp secret set 'mongodb.connection_string=mongodb://admin:password@host.docker.internal:27017/' && \
	docker mcp secret set 'resend.api_key=$(RESEND_API_KEY)' && \
	docker mcp secret set 'mysql.password=app_password' && \
	docker mcp secret set 'brave.api_key=$(BRAVE_API_KEY)' 

dev-start:
	MCPGATEWAY_ENDPOINT=http://localhost:8811/sse \
	MODEL_RUNNER_URL=http://localhost:12434/engines/llama.cpp/v1 \
	MODEL_RUNNER_MODEL=ai/qwen3:14B-Q6_K \
	uv run adk web --port 8080 --log_level DEBUG
