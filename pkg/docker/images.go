package docker

import (
	"context"
	"fmt"
	"io"
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
	errs, ctx := errgroup.WithContext(ctx)
	errs.SetLimit(runtime.NumCPU())

	for _, name := range names {
		errs.Go(func() error {
			return c.PullImage(ctx, name)
		})
	}

	return errs.Wait()
}

func (c *Client) PullImage(ctx context.Context, name string) error {
	response, err := c.client.ImagePull(ctx, name, image.PullOptions{})
	if err != nil {
		return fmt.Errorf("pulling docker image %s: %w", name, err)
	}

	if _, err := io.Copy(io.Discard, response); err != nil {
		return fmt.Errorf("pulling docker image %s: %w", name, err)
	}

	return nil
}
