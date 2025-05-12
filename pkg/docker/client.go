package docker

import (
	"context"

	"github.com/docker/cli/cli/command"
	"github.com/docker/cli/cli/context/docker"
	"github.com/docker/cli/cli/flags"
	"github.com/docker/docker/client"
)

type Client struct {
	client *client.Client
}

func NewClient(ctx context.Context) (*Client, error) {
	cli, err := command.NewDockerCli(command.WithBaseContext(ctx))
	if err != nil {
		return nil, err
	}

	if err := cli.Initialize(flags.NewClientOptions()); err != nil {
		return nil, err
	}

	currentContext := cli.CurrentContext()
	host := "/var/run/docker.sock"

	contexts, _ := cli.ContextStore().List()
	for _, c := range contexts {
		if c.Name == currentContext {
			dockerEndpoint, err := docker.EndpointFromContext(c)
			if err != nil {
				return nil, err
			}
			host = dockerEndpoint.Host
		}
	}

	c, err := client.NewClientWithOpts(client.FromEnv, client.WithHost(host), client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, err
	}

	return &Client{
		client: c,
	}, nil
}
