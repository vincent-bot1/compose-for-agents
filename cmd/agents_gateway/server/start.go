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
func startMCPClient(ctx context.Context, mcpImage string, pull bool, config string) (*mcpclient.Client, error) {
	serversByName, err := servers.List()
	if err != nil {
		return nil, fmt.Errorf("listing servers: %w", err)
	}

	var runConfig servers.Run
	if server, ok := serversByName[mcpImage]; ok {
		runConfig = server.Run
	}

	args := []string{"--security-opt", "no-new-privileges"}
	if runConfig.Workdir != "" {
		args = append(args, "--workdir", runConfig.Workdir)
	}
	// TODO: runConfig.Env
	// TODO: runConfig.Volumes
	// TODO: reeplace placeholders in runConfig.Command

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

	client := mcpclient.NewClientArgs(mcpImage, pull, env, args, runConfig.Command)
	if err := client.Start(ctx); err != nil {
		return nil, fmt.Errorf("failed to start server %s: %w", mcpImage, err)
	}

	return client, nil
}
