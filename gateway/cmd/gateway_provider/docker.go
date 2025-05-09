package main

import (
	"context"
	"fmt"
	"sync"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

var dockerClient = sync.OnceValue(func() *client.Client {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}
	return cli
})

func Exists(ctx context.Context, containerID string) (bool, container.InspectResponse, error) {
	response, err := dockerClient().ContainerInspect(ctx, containerID)
	if client.IsErrNotFound(err) {
		return false, response, nil
	}
	return err == nil, response, err
}

func RemoveContainer(ctx context.Context, containerID string, force bool) error {
	return dockerClient().ContainerRemove(ctx, containerID, container.RemoveOptions{
		Force: force,
	})
}

func StartContainer(ctx context.Context, containerID string, containerConfig container.Config, hostConfig container.HostConfig) error {
	cli := dockerClient()

	resp, err := cli.ContainerCreate(ctx, &containerConfig, &hostConfig, nil, nil, containerID)
	if err != nil {
		return fmt.Errorf("creating container: %w", err)
	}

	if err := cli.ContainerStart(ctx, resp.ID, container.StartOptions{}); err != nil {
		return fmt.Errorf("starting container: %w", err)
	}

	return nil
}
