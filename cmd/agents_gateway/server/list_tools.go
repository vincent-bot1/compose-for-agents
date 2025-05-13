package server

import (
	"context"
	"fmt"
	"os"
	"runtime"
	"strings"
	"sync"

	"github.com/docker/compose-agents-demo/pkg/catalog"
	"github.com/docker/compose-agents-demo/pkg/config"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"golang.org/x/sync/errgroup"
)

func (g *Gateway) listTools(ctx context.Context, mcpCatalog catalog.Catalog, registryConfig config.Registry, serverNames []string) ([]server.ServerTool, error) {
	var serverTools []server.ServerTool
	var serverToolsLock sync.Mutex

	errs, ctx := errgroup.WithContext(ctx)
	errs.SetLimit(runtime.NumCPU())
	for _, serverName := range serverNames {
		serverConfig, tools, found := mcpCatalog.Find(serverName)

		switch {
		case !found:
			fmt.Fprintln(os.Stderr, "MCP server not found:", serverName)

		case serverConfig != nil:
			serverName := serverName
			errs.Go(func() error {
				client, err := g.startMCPClient(ctx, *serverConfig, registryConfig)
				if err != nil {
					fmt.Fprintln(os.Stderr, "Can't start MCP server:", err)
					return nil
				}

				tools, err := client.ListTools(ctx)
				client.Close() // Close early
				if err != nil {
					fmt.Fprintln(os.Stderr, "Can't list tools:", err)
					return nil
				}

				for _, tool := range tools {
					if !isToolEnabled(serverName, serverConfig.Image, tool.Name, g.ToolsNames) {
						continue
					}

					serverTool := server.ServerTool{
						Tool:    tool,
						Handler: g.mcpServerHandler(*serverConfig, registryConfig, tool.Name),
					}

					serverToolsLock.Lock()
					serverTools = append(serverTools, serverTool)
					serverToolsLock.Unlock()
				}

				return nil
			})

		case tools != nil:
			for _, tool := range *tools {
				if !isToolEnabled(serverName, "", tool.Name, g.ToolsNames) {
					continue
				}

				mcpTool := mcp.Tool{
					Name:        tool.Name,
					Description: tool.Description,
				}
				if len(tool.Parameters.Properties) > 0 {
					mcpTool.InputSchema.Type = tool.Parameters.Type
					mcpTool.InputSchema.Properties = tool.Parameters.Properties.ToMap()
					mcpTool.InputSchema.Required = tool.Parameters.Required
				} else {
					mcpTool.InputSchema.Type = "object"
				}

				serverTool := server.ServerTool{
					Tool:    mcpTool,
					Handler: mcpToolHandler(tool),
				}

				serverToolsLock.Lock()
				serverTools = append(serverTools, serverTool)
				serverToolsLock.Unlock()
			}
		}
	}

	return serverTools, errs.Wait()
}

func isToolEnabled(serverName, serverImage, toolName string, enabledTools []string) bool {
	if len(enabledTools) == 0 {
		return true
	}

	for _, enabled := range enabledTools {
		if strings.EqualFold(enabled, toolName) ||
			strings.EqualFold(enabled, serverName+":"+toolName) ||
			strings.EqualFold(enabled, serverName+":*") ||
			strings.EqualFold(enabled, "*") {
			return true
		}
	}

	if serverImage != "" {
		for _, enabled := range enabledTools {
			if strings.EqualFold(enabled, serverImage+":"+toolName) ||
				strings.EqualFold(enabled, serverImage+":*") {
				return true
			}
		}
	}

	return false
}
