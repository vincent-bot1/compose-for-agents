package gateway

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	responses "github.com/docker/compose-agents-demo/pkg/mcp"
	"github.com/docker/compose-agents-demo/pkg/secretsscan"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func callbacks(logCalls, scanSecrets bool) server.ToolHandlerMiddleware {
	return func(next server.ToolHandlerFunc) server.ToolHandlerFunc {
		return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			start := time.Now()
			tool := request.Params.Name
			arguments := argumentsToString(request.Params.Arguments)

			if logCalls {
				fmt.Printf("- Calling tool %s with arguments: %s\n", tool, arguments)
			}

			if scanSecrets {
				fmt.Printf("  - Scanning tool call arguments for secrets...\n")
				if secretsscan.ContainsSecrets(arguments) {
					return nil, fmt.Errorf("a secret is being passed to tool %s", tool)
				}
				fmt.Printf("  > No secret found in arguments.\n")
			}

			result, err := next(ctx, request)
			if err != nil {
				return result, err
			}

			if scanSecrets {
				fmt.Printf("  - Scanning tool call response for secrets...\n")

				var contents string
				for _, content := range result.Content {
					switch c := content.(type) {
					case mcp.TextContent:
						contents += c.Text
					case *mcp.TextContent:
						contents += c.Text
					}
				}

				if secretsscan.ContainsSecrets(contents) {
					return responses.ToolError(fmt.Sprintf("a secret is being returned by the %s tool", tool)), nil
					// return nil, fmt.Errorf("a secret is being returned by the %s tool", tool)
				}
				fmt.Printf("  > No secret found in response.\n")
			}

			if logCalls {
				fmt.Printf("> Calling tool %s took: %s\n", tool, time.Since(start))
			}

			return result, nil
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
