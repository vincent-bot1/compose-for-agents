package main

import (
	"context"
	"crypto/sha256"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/spf13/cobra"

	"github.com/docker/compose-agents-demo/pkg/compose"
	"github.com/docker/compose-agents-demo/pkg/docker"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/go-connections/nat"
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

			return nil
		},
	}
}

func startGateway(ctx context.Context, serviceName string, flags Flags) error {
	const (
		agentsYaml = "/agents.yaml"
	)
	client, err := docker.NewClient(ctx)
	if err != nil {
		return err
	}

	cmd := []string{"/agents.yaml"}

	containerID := flags.ContainerName(serviceName)
	exists, inspect, err := client.Exists(ctx, containerID)
	if err != nil {
		return err
	}

	configHash := fmt.Sprintf("%x", sha256.Sum256([]byte(strings.Join(cmd, ", "))))
	if exists {
		if inspect.State.Running && inspect.Config.Labels[compose.LabelNames.ConfigHash] == configHash {
			return nil
		}
		if err := client.RemoveContainer(ctx, containerID, true); err != nil {
			return err
		}
	}

	agentsYamlSource := flags.Config
	if !filepath.IsAbs(agentsYamlSource) {
		abs, err := filepath.Abs(agentsYamlSource)
		if err != nil {
			return fmt.Errorf("determining absolute path for config: %w", err)
		}
		agentsYamlSource = abs
	}

	var portBindings nat.PortMap

	if flags.APIPort != "" {
		host, port, err := net.SplitHostPort(flags.APIPort)
		if err != nil {
			host = "127.0.0.1"
			port = flags.APIPort
		}
		portNum, err := strconv.Atoi(port)
		if err != nil {
			return fmt.Errorf("invalid API port number: %w", err)
		}
		portBindings = nat.PortMap{
			"7777/tcp": []nat.PortBinding{
				{
					HostIP:   host,
					HostPort: strconv.Itoa(portNum),
				},
			},
		}
	}

	return client.StartContainer(ctx, containerID, container.Config{
		Image: "demo/agents",
		Cmd:   cmd,
		Env:   append(os.Environ(), "OPENAI_API_KEY="+flags.OpenAIAPIKey),
		Labels: map[string]string{
			compose.LabelNames.Project:         flags.Project,
			compose.LabelNames.Service:         serviceName,
			compose.LabelNames.OneOff:          "False",
			compose.LabelNames.ContainerNumber: "1",
			compose.LabelNames.ConfigHash:      configHash,
		},
	}, container.HostConfig{
		PortBindings: portBindings,
		NetworkMode:  container.NetworkMode(flags.NetworkName()),
		Init:         &trueValue,
		Mounts: []mount.Mount{
			{
				Type:   mount.TypeBind,
				Source: agentsYamlSource,
				Target: agentsYaml,
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
