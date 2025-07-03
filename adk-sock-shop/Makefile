build-image:
	docker buildx build --builder hydrobuild --platform linux/amd64,linux/arm64 --tag jimclark106/supplier_agent:latest --push --file Dockerfile .

gateway-secrets:
	docker mcp secret set 'mongodb.connection_string=mongodb://admin:password@host.docker.internal:27017/' && \
	docker mcp secret set 'resend.api_key=$(RESEND_API_KEY)' && \
	docker mcp secret set 'brave.api_key=$(BRAVE_API_KEY)' && \
	docker mcp secret export brave resend mongodb > ./.mcp.env

api-server-start:
	MCPGATEWAY_ENDPOINT=http://localhost:8811/sse \
	MODEL_RUNNER_URL=http://localhost:12434/engines/llama.cpp/v1 \
	MODEL_RUNNER_MODEL=ai/qwen3:14B-Q6_K \
	uv run adk api_server --port 8080 --log_level DEBUG

ui-start:
	uv run streamlit run apps/vendor_app.py --server.port 3000

compose-local:
	docker compose up mongodb mcp-gateway
