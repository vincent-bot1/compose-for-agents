package server

import (
	"context"
	"fmt"
	"os"

	mcpclient "github.com/docker/compose-agents-demo/cmd/agents_gateway/mcp"
	"github.com/docker/compose-agents-demo/pkg/catalog"
	"github.com/docker/compose-agents-demo/pkg/config"
	"github.com/docker/compose-agents-demo/pkg/docker"
	"github.com/docker/compose-agents-demo/pkg/eval"
)

func (g *Gateway) startMCPClient(ctx context.Context, server catalog.Server, registryConfig config.Registry) (*mcpclient.Client, error) {
	args := []string{"--security-opt", "no-new-privileges"}
	if server.Run.Workdir != "" {
		args = append(args, "--workdir", server.Run.Workdir)
	}

	var env []string
	for _, s := range server.Config.Secrets {
		if g.Standalone {
			value, err := docker.SecretValue(ctx, s.Id)
			if err != nil {
				return nil, fmt.Errorf("getting secret %s: %w", s.Name, err)
			}

			env = append(env, fmt.Sprintf("%s=%s", s.Name, value))
		}

		args = append(args, "-e", s.Name)
	}

	configuration := registryConfig.Servers[server.Name]
	for _, e := range server.Config.Env {
		value := eval.Expression(e.Expression, configuration.Config)
		env = append(env, fmt.Sprintf("%s=%s", e.Name, value))
	}

	for _, mount := range eval.Expressions(server.Run.Volumes, configuration.Config) {
		args = append(args, "-v", mount)
	}

	command := eval.Expressions(server.Run.Command, configuration.Config)

	fmt.Fprintln(os.Stderr, "Starting server", server.Image, "with args", args, "and command", command)

	client := mcpclient.NewClientArgs(server.Image, false, env, args, command)
	if err := client.Start(ctx); err != nil {
		return nil, fmt.Errorf("failed to start server %s: %w", server.Image, err)
	}

	return client, nil
}
