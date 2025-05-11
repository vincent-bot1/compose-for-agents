package server

import (
	"context"
	"fmt"
	"net"

	"github.com/docker/compose-agents-demo/pkg/catalog"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func Run(ctx context.Context, serverNames, tools string, logCalls, scanSecrets bool) error {
	// List as early as possible to not lose client connections
	var lc net.ListenConfig
	ln, err := lc.Listen(ctx, "tcp", ":8811")
	if err != nil {
		return err
	}

	// Read the MCP catalog
	mcpCatalog, err := catalog.Get()
	if err != nil {
		return fmt.Errorf("listing catalog: %w", err)
	}

	toolCallbacks := callbacks(logCalls, scanSecrets)

	serverTools, err := listTools(ctx, serverNames, mcpCatalog.Servers, tools)
	if err != nil {
		return fmt.Errorf("listing tools: %w", err)
	}

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

				mcpServer := server.NewMCPServer("Docker AI MCP Gateway", "1.0.1", server.WithToolHandlerMiddleware(toolCallbacks))
				mcpServer.SetTools(serverTools...)
				stdioServer := server.NewStdioServer(mcpServer)

				if err := stdioServer.Listen(ctx, conn, conn); err != nil {
					fmt.Printf("Error listening: %v\n", err)
				}
			}()
		}
	}
}

func mcpServerHandler(server catalog.Server, tool mcp.Tool) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		client, err := startMCPClient(ctx, server, false)
		if err != nil {
			return nil, err
		}
		defer client.Close()

		return client.CallTool(ctx, tool.Name, request.Params.Arguments)
	}
}
