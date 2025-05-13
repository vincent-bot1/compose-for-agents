package server

import (
	"context"
	"fmt"

	mcpclient "github.com/docker/compose-agents-demo/cmd/agents_gateway/mcp"
	"github.com/docker/compose-agents-demo/pkg/catalog"
	"github.com/docker/compose-agents-demo/pkg/docker"
	"github.com/docker/compose-agents-demo/pkg/eval"
)

func (g *Gateway) startMCPClient(ctx context.Context, server catalog.Server, serverConfig map[string]any) (*mcpclient.Client, error) {
	image := server.Image
	command := eval.Expressions(server.Run.Command, serverConfig)
	args, env, err := g.argsAndEnv(ctx, server, serverConfig)
	if err != nil {
		return nil, err
	}

	if len(command) == 0 {
		log("Starting server", image, "with args", args)
	} else {
		log("Starting server", image, "with args", args, "and command", command)
	}

	client := mcpclient.NewClientArgs(image, false, env, args, command)
	if err := client.Start(ctx); err != nil {
		return nil, fmt.Errorf("failed to start server %s: %w", image, err)
	}

	return client, nil
}

func (g *Gateway) argsAndEnv(ctx context.Context, serverSpec catalog.Server, serverConfig map[string]any) ([]string, []string, error) {
	args := []string{"--security-opt", "no-new-privileges", "--cpus", "1", "--memory", "2Gb"}

	var env []string
	for _, s := range serverSpec.Config.Secrets {
		args = append(args, "-e", s.Name)

		if g.Standalone {
			value, err := docker.SecretValue(ctx, s.Id)
			if err != nil {
				return nil, nil, fmt.Errorf("getting secret %s: %w", s.Name, err)
			}

			env = append(env, fmt.Sprintf("%s=%s", s.Name, value))
		}
	}

	for _, e := range serverSpec.Config.Env {
		value := eval.Expression(e.Expression, serverConfig)
		env = append(env, fmt.Sprintf("%s=%s", e.Name, value))
	}

	for _, mount := range eval.Expressions(serverSpec.Run.Volumes, serverConfig) {
		args = append(args, "-v", mount)
	}

	return args, env, nil
}
