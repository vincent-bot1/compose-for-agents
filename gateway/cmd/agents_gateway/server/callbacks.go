package server

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/docker/gateway/cmd/agents_gateway/secrets"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func callbacks(logCalls, scanSecrets bool) server.ToolHandlerMiddleware {
	return func(next server.ToolHandlerFunc) server.ToolHandlerFunc {
		return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			tool := request.Params.Name
			arguments := argumentsToString(request.Params.Arguments)

			if scanSecrets {
				fmt.Printf("Scanning tool call arguments for secrets...\n")
				if secrets.ContainsSecrets(arguments) {
					return nil, fmt.Errorf("a secret is being passed to tool %s", tool)
				}
			}

			if logCalls {
				fmt.Printf("Calling tool %s with arguments: %s\n", tool, arguments)
			}

			return next(ctx, request)
		}
	}
}

func argumentsToString(args map[string]any) string {
	buf, err := json.Marshal(args)
	if err != nil {
		return fmt.Sprintf("%v", args)
	}

	return string(buf)
}
