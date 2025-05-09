package main

import (
	"context"
	"crypto/sha256"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
)

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
	cmd := []string{
		"--servers=" + flags.Servers,
		"--config=" + flags.Config,
		"--tools=" + flags.Tools,
		"--log_calls=" + boolToString(flags.LogCallsEnabled()),
		"--scan_secrets=" + boolToString(flags.ScanSecretsEnabled()),
	}

	containerID := flags.ContainerName(serviceName)
	exists, inspect, err := Exists(ctx, containerID)
	if err != nil {
		return err
	}

	configHash := fmt.Sprintf("%x", sha256.Sum256([]byte(strings.Join(cmd, ", "))))
	if exists {
		if inspect.State.Running && inspect.Config.Labels[labelNames.ConfigHash] == configHash {
			return nil
		}
		if err := RemoveContainer(ctx, containerID, true); err != nil {
			return err
		}
	}

	return StartContainer(ctx, containerID, container.Config{
		Image: flags.Image,
		Cmd:   cmd,
		Env: []string{
			// TEMP for github MCP server
			"GITHUB_TOKEN=" + os.Getenv("GITHUB_TOKEN"),
		},
		Labels: map[string]string{
			labelNames.Project:         flags.Project,
			labelNames.Service:         serviceName,
			labelNames.OneOff:          "False",
			labelNames.ContainerNumber: "1",
			labelNames.ConfigHash:      configHash,
		},
	}, container.HostConfig{
		NetworkMode: container.NetworkMode(flags.NetworkName()),
		Init:        &trueValue,
		Mounts: []mount.Mount{
			{
				Type:   mount.TypeBind,
				Source: "/var/run/docker.sock",
				Target: "/var/run/docker.sock",
			},
		},
	})
}

var trueValue = true

func boolToString(b bool) string {
	if b {
		return "true"
	}
	return "false"
}
