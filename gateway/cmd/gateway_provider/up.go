package main

import (
	"context"
	"crypto/sha256"
	"fmt"
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
	project := flags.Project
	containerID := flags.ContainerName(serviceName)
	network := flags.NetworkName()
	tools := flags.Tools
	logCalls := flags.LogCallsEnabled()
	scanSecrets := flags.ScanSecretsEnabled()

	cmd := []string{
		"--tools=" + tools,
		"--log_calls=" + boolToString(logCalls),
		"--scan_secrets=" + boolToString(scanSecrets),
	}

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
		Labels: map[string]string{
			labelNames.Project:         project,
			labelNames.Service:         serviceName,
			labelNames.OneOff:          "False",
			labelNames.ContainerNumber: "1",
			labelNames.ConfigHash:      configHash,
		},
	}, container.HostConfig{
		NetworkMode: container.NetworkMode(network),
		Mounts: []mount.Mount{
			{
				Type:   mount.TypeBind,
				Source: "/var/run/docker.sock",
				Target: "/var/run/docker.sock",
			},
		},
	})
}

func boolToString(b bool) string {
	if b {
		return "true"
	}
	return "false"
}
