package server

import (
	"context"
	"fmt"

	mcpclient "github.com/docker/compose-agents-demo/cmd/agents_gateway/mcp"
	"github.com/docker/compose-agents-demo/pkg/catalog"
)

func startMCPClient(ctx context.Context, server catalog.Server, pull bool) (*mcpclient.Client, error) {
	args := []string{"--security-opt", "no-new-privileges"}
	if server.Run.Workdir != "" {
		args = append(args, "--workdir", server.Run.Workdir)
	}
	// TODO: runConfig.Env
	// TODO: runConfig.Volumes
	// TODO: reeplace placeholders in runConfig.Command

	var env []string
	for _, secret := range server.Config.Secrets {
		args = append(args, "-e", secret.Name)
	}

	client := mcpclient.NewClientArgs(server.Image, pull, env, args, server.Run.Command)
	if err := client.Start(ctx); err != nil {
		return nil, fmt.Errorf("failed to start server %s: %w", server.Image, err)
	}

	return client, nil
}
