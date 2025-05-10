package server

import (
	"context"
	"fmt"
	"net"

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
			conn, err := acceptWithContext(ctx, ln)
			if err != nil {
				if ctx.Err() != nil {
					return ctx.Err()
				}
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
