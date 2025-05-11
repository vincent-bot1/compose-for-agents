package server

import (
	"context"
	"fmt"
	"os/exec"

	"github.com/docker/compose-agents-demo/pkg/catalog"
	"github.com/docker/compose-agents-demo/pkg/eval"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func mcpToolHandler(tool catalog.Tool) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		command, err := eval.Expressions(tool.Container.Command, request.Params.Arguments)
		if err != nil {
			return nil, fmt.Errorf("replacing placeholders: %w", err)
		}

		args := []string{"run", "--rm", "-i", "--init", "--security-opt", "no-new-privileges"}
		for _, v := range tool.Container.Volumes {
			args = append(args, "-v", v)
		}
		args = append(args, tool.Container.Image)
		args = append(args, command...)

		cmd := exec.CommandContext(ctx, "docker", args...)
		out, err := cmd.CombinedOutput()
		if err != nil {
			return toolError(string(out)), nil
		}

		return toolResult(string(out)), nil
	}
}
