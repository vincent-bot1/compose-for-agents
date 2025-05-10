package server

import (
	"context"
	"fmt"
	"net"
	"os"
	"strings"

	mcpclient "github.com/docker/compose-agents-demo/gateway/cmd/agents_gateway/mcp"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func Run(ctx context.Context, servers, config, tools string, logCalls, scanSecrets bool) error {
	// List as early as possible to not lose client connections
	var lc net.ListenConfig
	ln, err := lc.Listen(ctx, "tcp", ":8811")
	if err != nil {
		return err
	}

	serverTools, err := listTools(ctx, servers, tools, config)
	if err != nil {
		return fmt.Errorf("listing tools: %w", err)
	}

	mcpServer := server.NewMCPServer(
		"Docker AI MCP Gateway",
		"1.0.1",
		server.WithToolHandlerMiddleware(callbacks(logCalls, scanSecrets)),
	)
	mcpServer.SetTools(serverTools...)

	stdioServer := server.NewStdioServer(mcpServer)
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			conn, err := ln.Accept()
			if err != nil {
				fmt.Printf("Error accepting the connection: %v\n", err)
				continue
			}

			go func() {
				defer conn.Close()
				if err := stdioServer.Listen(ctx, conn, conn); err != nil {
					fmt.Printf("Error listening: %v\n", err)
				}
			}()
		}
	}
}

func mcpServerHandler(mcpImage string, tool mcp.Tool, config string) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		client, err := startMCPClient(ctx, mcpImage, false, config)
		if err != nil {
			return nil, err
		}
		defer client.Close()

		return client.CallTool(ctx, tool.Name, request.Params.Arguments)
	}
}

// config: mcp/github-mcp-server.GITHUB_PERSONAL_ACCESS_TOKEN=$GITHUB_TOKEN
func startMCPClient(ctx context.Context, mcpImage string, pull bool, config string) (*mcpclient.Client, error) {
	var args []string
	var env []string

	for _, cfg := range parseConfig(config) {
		prefix := mcpImage + "."
		if !strings.HasPrefix(cfg, prefix) {
			continue
		}

		mapping := strings.TrimPrefix(cfg, prefix)
		parts := strings.SplitN(mapping, "=", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid config format: %s", cfg)
		}

		if strings.HasPrefix(parts[1], "$") {
			env = append(args, parts[0]+"="+os.Getenv(parts[1][1:]))
		} else {
			env = append(args, parts[0]+"="+parts[1])
		}
		args = append(args, "-e", parts[0])
	}

	client := mcpclient.NewClientArgs(mcpImage, pull, env, args, nil)
	if err := client.Start(ctx); err != nil {
		return nil, fmt.Errorf("failed to start server %s: %w", mcpImage, err)
	}

	return client, nil
}
