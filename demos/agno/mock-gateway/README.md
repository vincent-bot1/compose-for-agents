# Mock Gateway

A Go-based mock service that simulates MCP (Model Context Protocol) tool responses for testing purposes.

## Features

- **Dynamic Mock Loading**: Automatically loads mock responses from JSON files
- **RESTful API**: Provides endpoints for health checks and mock management
- **Flexible Configuration**: Easy to add new mocks by adding JSON files
- **Docker Support**: Containerized for easy deployment

## API Endpoints

### Health Check
```
GET /health
```
Returns the service status and number of loaded mocks.

### List Mocks
```
GET /mocks
```
Returns a list of all loaded mock configurations.

### MCP Tool Calls
```
POST /mcp/{tool}/{method}
```
Handles MCP tool calls and returns configured mock responses.

## Mock Configuration

Mocks are defined as JSON files in the `mocks/` directory. Each file should contain:

```json
{
  "tool": "tool-name",
  "method": "method-name", 
  "response": {
    // Your mock response data here
  }
}
```

### Example Mock Files

- `github-list-issues.json`: Mocks GitHub issue listing
- `notion-search.json`: Mocks Notion page search
- `notion-create-page.json`: Mocks Notion page creation
- `fetch-content.json`: Mocks content fetching

## Environment Variables

- `PORT`: Server port (default: 8080)
- `MOCKS_DIR`: Directory containing mock JSON files (default: ./mocks)

## Usage

### Local Development
```bash
go run main.go
```

### Docker
```bash
docker build -t mock-gateway .
docker run -p 8080:8080 -v $(pwd)/mocks:/root/mocks mock-gateway
```

### Testing
```bash
# Health check
curl http://localhost:8080/health

# List mocks
curl http://localhost:8080/mocks

# Test a mock call
curl -X POST http://localhost:8080/mcp/github-mcp-server/list_issues \
  -H "Content-Type: application/json" \
  -d '{"repository": "example/repo"}'
```

## Integration with Compose

The mock-gateway can be integrated into your docker-compose setup to replace the real MCP gateway during testing.

## Adding New Mocks

1. Create a new JSON file in the `mocks/` directory
2. Define the tool, method, and response structure
3. Restart the service to load the new mock
4. The mock will be automatically available at `/mcp/{tool}/{method}`
