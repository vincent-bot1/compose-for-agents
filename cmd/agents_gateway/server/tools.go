package server

import (
	"context"
	"fmt"

	"github.com/docker/compose-agents-demo/pkg/catalog"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func mcpToolHandler(tool catalog.Tool) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return nil, fmt.Errorf("tool %s not implemented", tool.Name)
	}
}
