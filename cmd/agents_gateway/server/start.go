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

	var env []string
	for _, s := range server.Config.Secrets {
		args = append(args, "-e", s.Name)
	}
	for _, e := range server.Config.Env {
		args = append(args, "-e", e.Name)
	}
	// TODO: replace placeholders in server.Run.Volumes
	for _, v := range server.Run.Volumes {
		args = append(args, "-v", v)
	}

	// TODO: replace placeholders in server.Run.Command
	client := mcpclient.NewClientArgs(server.Image, pull, env, args, server.Run.Command)
	if err := client.Start(ctx); err != nil {
		return nil, fmt.Errorf("failed to start server %s: %w", server.Image, err)
	}

	return client, nil
}
