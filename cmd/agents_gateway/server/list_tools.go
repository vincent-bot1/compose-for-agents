package server

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"github.com/docker/compose-agents-demo/pkg/catalog"
	"github.com/docker/compose-agents-demo/pkg/config"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"golang.org/x/sync/errgroup"
)

func listTools(ctx context.Context, mcpCatalog catalog.Catalog, registryConfig config.Registry, serverNames []string, toolNames []string) ([]server.ServerTool, error) {
	var serverTools []server.ServerTool
	var serverToolsLock sync.Mutex

	errs, ctx := errgroup.WithContext(ctx)
	for _, serverName := range serverNames {
		// Is it an MCP Server?
		serverConfig, ok := mcpCatalog.Servers[serverName]
		if !ok {
			// Is it a tool group?
			tools, ok := mcpCatalog.Tools[serverName]
			if !ok {
				fmt.Println("MCP server not found:", serverName)
				continue
			}

			for _, tool := range tools {
				if !isToolEnabled(serverName, tool.Name, toolNames) {
					continue
				}

				serverTool := server.ServerTool{
					Tool: mcp.Tool{
						Name:        tool.Name,
						Description: tool.Description,
						InputSchema: mcp.ToolInputSchema{
							Type:       tool.Parameters.Type,
							Properties: tool.Parameters.Properties.ToMap(),
							Required:   tool.Parameters.Required,
						},
					},
					Handler: mcpToolHandler(tool),
				}

				serverToolsLock.Lock()
				serverTools = append(serverTools, serverTool)
				serverToolsLock.Unlock()
			}

			continue
		}

		serverName := serverName
		errs.Go(func() error {
			client, err := startMCPClient(ctx, serverConfig, registryConfig)
			if err != nil {
				fmt.Println("Can't start MCP server:", err)
				return nil
			}

			tools, err := client.ListTools(ctx)
			client.Close() // Close early
			if err != nil {
				fmt.Println("Can't list tools:", err)
				return nil
			}

			for _, tool := range tools {
				if !isToolEnabled(serverName, tool.Name, toolNames) {
					continue
				}

				serverTool := server.ServerTool{
					Tool:    tool,
					Handler: mcpServerHandler(serverConfig, registryConfig, tool),
				}

				serverToolsLock.Lock()
				serverTools = append(serverTools, serverTool)
				serverToolsLock.Unlock()
			}

			return nil
		})
	}

	return serverTools, errs.Wait()
}

func isToolEnabled(serverName string, toolName string, toolsNames []string) bool {
	for _, enabled := range toolsNames {
		if strings.EqualFold(enabled, toolName) ||
			strings.EqualFold(enabled, serverName+":"+toolName) ||
			strings.EqualFold(enabled, "mcp/"+serverName+":"+toolName) ||
			strings.EqualFold(enabled, serverName+":*") ||
			strings.EqualFold(enabled, "mcp/"+serverName+":*") ||
			strings.EqualFold(enabled, "*") {
			return true
		}
	}

	return false
}
