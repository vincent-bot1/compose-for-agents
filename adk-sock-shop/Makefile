build-image:
	docker buildx build --builder hydrobuild --platform linux/amd64,linux/arm64 --tag jimclark106/supplier_agent:latest --push --file Dockerfile .

gateway-secrets:
	docker mcp secret set 'mongodb.connection_string=mongodb://admin:password@mongodb:27017/'

dev-start:
	MCPGATEWAY_ENDPOINT=http://localhost:8811/sse \
	MODEL_RUNNER_URL=http://localhost:12434/engines/llama.cpp/v1 \
	MODEL_RUNNER_MODEL=ai/gemma3:4B-Q4_0 \
	uv run adk web --port 8080 --log_level DEBUG
