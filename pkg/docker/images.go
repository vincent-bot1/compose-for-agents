package docker

import (
	"context"
	"fmt"
	"io"
	"os"
	"runtime"

	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/client"
	"golang.org/x/sync/errgroup"
)

func (c *Client) ImageExists(ctx context.Context, name string) (bool, error) {
	_, err := c.client.ContainerInspect(ctx, name)
	if client.IsErrNotFound(err) {
		return false, nil
	}

	return err == nil, err
}

func (c *Client) PullImages(ctx context.Context, names ...string) error {
	token, err := getJWT(ctx)
	if err != nil {
		return fmt.Errorf("getting auth token: %w", err)
	}

	errs, ctx := errgroup.WithContext(ctx)
	errs.SetLimit(runtime.NumCPU())

	for _, name := range names {
		errs.Go(func() error {
			return c.PullImage(ctx, name, token)
		})
	}

	return errs.Wait()
}

func (c *Client) PullImage(ctx context.Context, name, token string) error {
	if token == "" {
		var err error
		token, err = getJWT(ctx)
		if err != nil {
			return fmt.Errorf("getting auth token: %w", err)
		}
	}

	pullOptions := image.PullOptions{
		RegistryAuth: token,
	}
	response, err := c.client.ImagePull(ctx, name, pullOptions)
	if err != nil {
		return fmt.Errorf("pulling docker image %s: %w", name, err)
	}

	if _, err := io.Copy(io.Discard, response); err != nil {
		return fmt.Errorf("pulling docker image %s: %w", name, err)
	}

	return nil
}

func getJWT(ctx context.Context) (string, error) {
	if _, err := os.Stat("/run/host-services/backend.sock"); err != nil {
		return "", nil
	}

	var token string
	if err := get(ctx, httpClient(dialHostSideBackend), "/registry/token", &token); err != nil {
		return "", fmt.Errorf("getting auth token: %w", err)
	}
	return token, nil
}
