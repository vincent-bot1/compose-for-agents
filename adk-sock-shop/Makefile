build-adk-image:
	docker buildx build --builder hydrobuild --platform linux/amd64,linux/arm64 --tag jimclark106/supplier_agent:latest --push --file Dockerfile .

build-adk-ui-image:
	docker buildx build --builder hydrobuild --platform linux/amd64,linux/arm64 --tag jimclark106/supplier_agent_ui:latest --push --file Dockerfile.adk-ui .

gateway-secrets:
	docker mcp secret set 'mongodb.connection_string=mongodb://admin:password@mongodb:27017/' && \
	docker mcp secret set 'resend.api_key=$(RESEND_API_KEY)' && \
	docker mcp secret set 'brave.api_key=$(BRAVE_API_KEY)' && \
	docker mcp secret export brave resend mongodb > ./.mcp.env && \
	echo $(OPENAI_API_KEY) > ./secret.openai-api-key

api-server-start:
	MCPGATEWAY_ENDPOINT=http://localhost:8811/sse \
	MODEL_RUNNER_URL=http://localhost:12434/engines/llama.cpp/v1 \
	MODEL_RUNNER_MODEL=ai/qwen3:14B-Q6_K \
	uv run adk api_server --port 8000 --log_level DEBUG

ui-start:
	API_BASE_URL=http://localhost:8000 \
	uv run streamlit run apps/vendor_app.py --server.port 3000

dev-context:
	docker context use desktop-linux

cloud-context:
	docker context use docker-cloud

dev-compose-up:
	docker compose up mongodb mcp-gateway

compose-down:
	docker compose down

dev-up: dev-context dev-compose-up

dev-down: dev-context compose-down

compose-sockstore:
	docker compose up frontend edge-router catalogue catalogue-db carts carts-db orders orders-db shipping queue-number rabbitmq payment user user-db user-sim

compose-agent:
	docker compose up mongodb mcp-gateway adk adk-ui --build

cloud-compose-up:
	docker compose -f compose.yaml -f compose.offload.yaml up --build

offload-up: cloud-context cloud-compose-up

cloud-partial-down:
	docker compose -f compose.yaml -f compose.offload.yaml down mcp-gateway adk adk-ui

