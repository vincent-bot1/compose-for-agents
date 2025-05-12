package server

import (
	"context"
	"os/exec"

	"github.com/docker/compose-agents-demo/pkg/catalog"
	"github.com/docker/compose-agents-demo/pkg/config"
	"github.com/docker/compose-agents-demo/pkg/eval"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func (g *Gateway) mcpServerHandler(server catalog.Server, registryConfig config.Registry, tool string) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		client, err := g.startMCPClient(ctx, server, registryConfig)
		if err != nil {
			return nil, err
		}
		defer client.Close()

		return client.CallTool(ctx, tool, request.Params.Arguments)
	}
}

func mcpToolHandler(tool catalog.Tool) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := []string{"run", "--rm", "-i", "--init", "--security-opt", "no-new-privileges"}

		for _, v := range eval.Expressions(tool.Container.Volumes, request.Params.Arguments) {
			args = append(args, "-v", v)
		}
		args = append(args, tool.Container.Image)
		args = append(args, eval.Expressions(tool.Container.Command, request.Params.Arguments)...)

		cmd := exec.CommandContext(ctx, "docker", args...)
		out, err := cmd.CombinedOutput()
		if err != nil {
			return toolError(string(out)), nil
		}

		return toolResult(string(out)), nil
	}
}
