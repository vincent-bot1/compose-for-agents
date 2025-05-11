package server

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"github.com/docker/compose-agents-demo/pkg/catalog"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"golang.org/x/sync/errgroup"
)

func listTools(ctx context.Context, serverNames string, mcpCatalog catalog.Catalog, tools string) ([]server.ServerTool, error) {
	// Filter out tools
	toolNeeded := map[string]bool{}
	for tool := range strings.SplitSeq(tools, ",") {
		toolNeeded[strings.TrimSpace(tool)] = true
	}

	var serverTools []server.ServerTool
	var serverToolsLock sync.Mutex

	errs, ctx := errgroup.WithContext(ctx)
	for _, serverName := range parseServers(serverNames) {
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
				if _, ok := toolNeeded[tool.Name]; !ok {
					continue
				}

				serverTool := server.ServerTool{
					Tool: mcp.Tool{
						Name:        tool.Name,
						Description: tool.Description,
						InputSchema: mcp.ToolInputSchema{
							Type: tool.Parameters.Type,
							// Properties: tool.Parameters.Properties,
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

		errs.Go(func() error {
			client, err := startMCPClient(ctx, serverConfig, true)
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
				if _, ok := toolNeeded[tool.Name]; !ok {
					continue
				}

				serverTool := server.ServerTool{
					Tool:    tool,
					Handler: mcpServerHandler(serverConfig, tool),
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
