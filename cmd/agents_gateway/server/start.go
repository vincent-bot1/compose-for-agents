package server

import (
	"context"
	"fmt"
	"os"
	"strings"

	mcpclient "github.com/docker/compose-agents-demo/cmd/agents_gateway/mcp"
	"github.com/docker/compose-agents-demo/cmd/agents_gateway/servers"
)

// config: mcp/github-mcp-server.GITHUB_PERSONAL_ACCESS_TOKEN=$GITHUB_TOKEN
func startMCPClient(ctx context.Context, mcpServer string, pull bool, config string) (*mcpclient.Client, error) {
	serversByName, err := servers.List()
	if err != nil {
		return nil, fmt.Errorf("listing servers: %w", err)
	}

	server, ok := serversByName[mcpServer]
	if !ok {
		return nil, fmt.Errorf("server not found: %s", mcpServer)
	}

	args := []string{"--security-opt", "no-new-privileges"}
	if server.Run.Workdir != "" {
		args = append(args, "--workdir", server.Run.Workdir)
	}
	// TODO: runConfig.Env
	// TODO: runConfig.Volumes
	// TODO: reeplace placeholders in runConfig.Command

	var env []string
	for _, cfg := range parseConfig(config) {
		prefix := server.Image + "."
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

	client := mcpclient.NewClientArgs(server.Image, pull, env, args, server.Run.Command)
	if err := client.Start(ctx); err != nil {
		return nil, fmt.Errorf("failed to start server %s: %w", mcpServer, err)
	}

	return client, nil
}
