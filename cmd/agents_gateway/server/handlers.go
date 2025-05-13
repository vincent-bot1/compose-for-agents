package server

import (
	"context"
	"os/exec"

	"github.com/docker/compose-agents-demo/pkg/catalog"
	"github.com/docker/compose-agents-demo/pkg/eval"
	responses "github.com/docker/compose-agents-demo/pkg/mcp"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func (g *Gateway) mcpServerHandler(serverSpec catalog.Server, serverConfig map[string]any, tool string) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		client, err := g.startMCPClient(ctx, serverSpec, serverConfig)
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
			return responses.ToolError(string(out)), nil
		}

		return responses.ToolResult(string(out)), nil
	}
}
