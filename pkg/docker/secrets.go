package docker

import (
	"context"
)

type Secret struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

func SecretValue(ctx context.Context, id string) (string, error) {
	// Make sure to always talk to Docker Desktop directly in order to read the secrets used by the MCP Toolkit extension.
	return RunOnDockerDesktop(ctx, "-l", "x-secret:"+id+"=/secret", "busybox", "cat", "/secret")
}
