package server

import (
	"context"
	"fmt"
	"os/exec"
	"time"
)

func VerifySignatures(ctx context.Context, mcpImages []string) error {
	start := time.Now()
	log("Verifying docker images", mcpImages)

	args := []string{"verify"}
	args = append(args, mcpImages...)
	args = append(args, "--key", "https://raw.githubusercontent.com/docker/keyring/refs/heads/main/public/mcp/latest.pub")

	cmd := exec.CommandContext(ctx, "/usr/bin/cosign", args...)
	cmd.Env = []string{"COSIGN_REPOSITORY=mcp/signatures"}

	if out, err := cmd.CombinedOutput(); err != nil {
		log("Failed to verify docker images:", string(out))
		return fmt.Errorf("verifying images: %w", err)
	}

	log("Docker images verified in", time.Since(start))
	return nil
}
