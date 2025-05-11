package main

import (
	"context"
	"crypto/sha256"
	"fmt"
	"sort"
	"strings"

	"github.com/spf13/cobra"

	"github.com/docker/compose-agents-demo/pkg/catalog"
	"github.com/docker/compose-agents-demo/pkg/compose"
	"github.com/docker/compose-agents-demo/pkg/docker"
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
				compose.ErrorMessage("could not start the gateway", err)
			} else {
				compose.InfoMessage("started the gateway")
			}

			compose.Setenv("ENDPOINT", flags.ContainerName(serviceName)+":8811")
			return nil
		},
	}
}

func startGateway(ctx context.Context, serviceName string, flags Flags) error {
	client, err := docker.NewClient(ctx)
	if err != nil {
		return err
	}

	enabledServers, err := enabledMCPServers(ctx)
	if err != nil {
		return fmt.Errorf("reading configuration: %w", err)
	}

	var serverNames []string
	for serverName := range enabledServers {
		serverNames = append(serverNames, serverName)
	}
	sort.Strings(serverNames)

	serversByName, err := catalog.Get()
	if err != nil {
		return fmt.Errorf("listing servers: %w", err)
	}

	var env []string
	for _, serverName := range serverNames {
		server, ok := serversByName[serverName]
		if !ok {
			continue
		}
		configuration, ok := enabledServers[serverName]
		if !ok {
			continue
		}

		for _, s := range server.Config.Secrets {
			value, err := secretValue(ctx, s.Id)
			if err != nil {
				return fmt.Errorf("getting secret %s: %w", s.Name, err)
			}

			env = append(env, fmt.Sprintf("%s=%s", s.Name, value))
		}
		for _, e := range server.Config.Env {
			env = append(env, fmt.Sprintf("%s=%s", e.Name, evaluate(e.Expression, configuration.Config[e.Name])))
		}
	}

	cmd := []string{
		"--servers=" + strings.Join(serverNames, ","),
		"--tools=" + flags.Tools,
		"--log_calls=" + boolToString(flags.LogCallsEnabled()),
		"--scan_secrets=" + boolToString(flags.ScanSecretsEnabled()),
	}

	containerID := flags.ContainerName(serviceName)
	exists, inspect, err := client.Exists(ctx, containerID)
	if err != nil {
		return err
	}

	// Make sure to restart the gateway if the config changes.
	configStr := string(catalog.McpServersYAML) + ":" + strings.Join(cmd, ", ") + ":" + strings.Join(env, ", ") + ":" + fmt.Sprintf("%+v", enabledServers)
	configHash := fmt.Sprintf("%x", sha256.Sum256([]byte(configStr)))
	if exists {
		if inspect.State.Running && inspect.Config.Labels[compose.LabelNames.ConfigHash] == configHash {
			return nil
		}
		if err := client.RemoveContainer(ctx, containerID, true); err != nil {
			return err
		}
	}

	return client.StartContainer(ctx, containerID, container.Config{
		Image: flags.Image,
		Cmd:   cmd,
		Env:   env,
		Labels: map[string]string{
			compose.LabelNames.Project:         flags.Project,
			compose.LabelNames.Service:         serviceName,
			compose.LabelNames.OneOff:          "False",
			compose.LabelNames.ContainerNumber: "1",
			compose.LabelNames.ConfigHash:      configHash,
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
