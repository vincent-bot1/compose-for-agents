package main

import (
	"context"
	"crypto/sha256"
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/docker/docker/api/types/container"
)

const gatewayImage = "docker/agents_gateway"

func NewUpCmd(flags *Flags) *cobra.Command {
	return &cobra.Command{
		Use:   "up",
		Short: "called during compose up",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			serviceName := args[0]

			if err := startGateway(cmd.Context(), serviceName, *flags); err != nil {
				errorMessage("could not start the gateway", err)
			} else {
				infoMessage("started the gateway")
			}

			setenv("ENDPOINT", flags.ContainerName(serviceName)+":8811")
			return nil
		},
	}
}

func startGateway(ctx context.Context, serviceName string, flags Flags) error {
	project := flags.Project
	containerID := flags.ContainerName(serviceName)
	network := flags.NetworkName()
	tools := flags.Tools
	logCalls := strings.EqualFold(flags.LogCalls, "yes")
	scanSecrets := strings.EqualFold(flags.ScanSecrets, "yes")

	cmd := []string{
		"--tools=" + tools,
		"--log_calls=" + boolToString(logCalls),
		"--scan_secrets=" + boolToString(scanSecrets),
	}

	configHash := fmt.Sprintf("%x", sha256.Sum256([]byte(strings.Join(cmd, ", "))))

	exists, inspect, err := Exists(ctx, containerID)
	if err != nil {
		return err
	}

	if exists && inspect.State.Running && inspect.Config.Labels["com.docker.compose.config-hash"] == configHash {
		return nil
	}

	if exists {
		if err := RemoveContainer(ctx, containerID, true); err != nil {
			return err
		}
	}

	return StartContainer(ctx, containerID, container.Config{
		Image: gatewayImage,
		Cmd:   cmd,
		Labels: map[string]string{
			"com.docker.compose.project":          project,
			"com.docker.compose.service":          serviceName,
			"com.docker.compose.oneoff":           "False",
			"com.docker.compose.container-number": "1",
			"com.docker.compose.config-hash":      configHash,
		},
	}, container.HostConfig{
		NetworkMode: container.NetworkMode(network),
	})
}

func boolToString(b bool) string {
	if b {
		return "true"
	}
	return "false"
}
