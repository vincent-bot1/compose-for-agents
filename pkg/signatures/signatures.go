package signatures

import (
	"context"
	"os/exec"
)

func Verify(ctx context.Context, images []string) error {
	args := []string{"verify"}
	args = append(args, images...)
	args = append(args, "--key", "https://raw.githubusercontent.com/docker/keyring/refs/heads/main/public/mcp/latest.pub")

	// TODO(dga): Could we replace cosign with our code?
	cmd := exec.CommandContext(ctx, "/usr/bin/cosign", args...)
	cmd.Env = []string{"COSIGN_REPOSITORY=mcp/signatures"}
	return cmd.Run()
}
