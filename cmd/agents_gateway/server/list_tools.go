package server

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"github.com/mark3labs/mcp-go/server"
	"golang.org/x/sync/errgroup"
)

func listTools(ctx context.Context, servers, tools, config string) ([]server.ServerTool, error) {
	// Filter out tools
	toolNeeded := map[string]bool{}
	for tool := range strings.SplitSeq(tools, ",") {
		toolNeeded[strings.TrimSpace(tool)] = true
	}

	var serverTools []server.ServerTool
	var serverToolsLock sync.Mutex

	errs, ctx := errgroup.WithContext(ctx)
	for _, mcpServer := range parseServers(servers) {
		mcpServer := mcpServer

		errs.Go(func() error {
			client, err := startMCPClient(ctx, mcpServer, true, config)
			if err != nil {
				return err
			}

			tools, err := client.ListTools(ctx)
			client.Close() // Close early
			if err != nil {
				return fmt.Errorf("listing tools: %w", err)
			}

			for _, tool := range tools {
				if _, ok := toolNeeded[tool.Name]; !ok {
					continue
				}

				serverTool := server.ServerTool{
					Tool:    tool,
					Handler: mcpServerHandler(mcpServer, tool, config),
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
