package docker

import (
	"context"
	"fmt"

	"github.com/docker/cli/cli/command"
	"github.com/docker/cli/cli/context/docker"
	"github.com/docker/cli/cli/flags"
	"github.com/docker/docker/api/types/container"
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

func (c *Client) Exists(ctx context.Context, containerID string) (bool, container.InspectResponse, error) {
	response, err := c.client.ContainerInspect(ctx, containerID)
	if client.IsErrNotFound(err) {
		return false, response, nil
	}

	return err == nil, response, err
}

func (c *Client) RemoveContainer(ctx context.Context, containerID string, force bool) error {
	return c.client.ContainerRemove(ctx, containerID, container.RemoveOptions{
		Force: force,
	})
}

func (c *Client) StartContainer(ctx context.Context, containerID string, containerConfig container.Config, hostConfig container.HostConfig) error {
	resp, err := c.client.ContainerCreate(ctx, &containerConfig, &hostConfig, nil, nil, containerID)
	if err != nil {
		return fmt.Errorf("creating container: %w", err)
	}

	if err := c.client.ContainerStart(ctx, resp.ID, container.StartOptions{}); err != nil {
		return fmt.Errorf("starting container: %w", err)
	}

	return nil
}
