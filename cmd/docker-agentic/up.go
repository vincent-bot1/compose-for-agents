package main

import (
	"context"
	"crypto/sha256"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/docker/compose-agents-demo/pkg/compose"
	"github.com/docker/compose-agents-demo/pkg/docker"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/go-connections/nat"
)

const (
	agentsImageName = "demo/agents"
	uiImageName     = "demo/ui"
)

func NewUpCmd(flags *Flags) *cobra.Command {
	return &cobra.Command{
		Use:   "up",
		Short: "called during compose up",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			serviceName := args[0]

			client, err := docker.NewClient(cmd.Context())
			if err != nil {
				return err
			}

			if err := startAgents(cmd.Context(), client, serviceName, flags); err != nil {
				compose.ErrorMessage("could not start agents", err)
			} else {
				compose.InfoMessage("agents started")
			}

			if flags.UIPort != "" {
				if err := startUI(cmd.Context(), client, serviceName, flags); err != nil {
					compose.ErrorMessage("could not start UI", err)
				} else {
					compose.InfoMessage("UI started")
				}
			}

			return nil
		},
	}
}

func startAgents(ctx context.Context, client *docker.Client, serviceName string, flags *Flags) error {
	const (
		agentsYaml = "/agents.yaml"
	)

	cmd := []string{"/agents.yaml"}

	containerID := flags.AgentsContainerName(serviceName)
	exists, inspect, err := client.Exists(ctx, containerID)
	if err != nil {
		return err
	}

	agentsData, err := os.ReadFile(flags.Config)
	if err != nil {
		return fmt.Errorf("reading agents.yaml: %w", err)
	}

	configHasher := sha256.New()
	if _, err := configHasher.Write(agentsData); err != nil {
		return fmt.Errorf("hashing agents.yaml: %w", err)
	}
	if _, err := configHasher.Write([]byte(flags.OpenAIAPIKey)); err != nil {
		return fmt.Errorf("hashing OpenAI API key: %w", err)
	}

	configHash := fmt.Sprintf("%x", configHasher.Sum(nil))
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
		Image: agentsImageName,
		Cmd:   cmd,
		// XXX: This needs the full environment due to MCPGATEWAY_ENDPOINT, one we figure out
		// a solution for this, pass only the required environment variables.
		Env: append(os.Environ(), []string{"OPENAI_API_KEY=" + flags.OpenAIAPIKey}...),
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

func startUI(ctx context.Context, client *docker.Client, serviceName string, flags *Flags) error {
	containerID := flags.UIContainerName(serviceName)
	exists, inspect, err := client.Exists(ctx, containerID)
	if err != nil {
		return err
	}

	if exists {
		if inspect.State.Running {
			return nil
		}
		if err := client.RemoveContainer(ctx, containerID, true); err != nil {
			return err
		}
	}

	var portBindings nat.PortMap

	var defaultEndpoint string
	if flags.APIPort != "" {
		_, port, err := net.SplitHostPort(flags.APIPort)
		if err != nil {
			port = flags.APIPort
		}
		defaultEndpoint = "http://localhost:" + port
	}

	if flags.UIPort != "" {
		host, port, err := net.SplitHostPort(flags.UIPort)
		if err != nil {
			host = "127.0.0.1"
			port = flags.UIPort
		}
		portNum, err := strconv.Atoi(port)
		if err != nil {
			return fmt.Errorf("invalid UI port number: %w", err)
		}
		portBindings = nat.PortMap{
			"3000/tcp": []nat.PortBinding{
				{
					HostIP:   host,
					HostPort: strconv.Itoa(portNum),
				},
			},
		}
	}

	return client.StartContainer(ctx, containerID, container.Config{
		Image: uiImageName,
		Env:   []string{"DEFAULT_ENDPOINT=" + defaultEndpoint},
		Labels: map[string]string{
			compose.LabelNames.Project:         flags.Project,
			compose.LabelNames.Service:         serviceName,
			compose.LabelNames.OneOff:          "False",
			compose.LabelNames.ContainerNumber: "1",
		},
	}, container.HostConfig{
		PortBindings: portBindings,
		NetworkMode:  container.NetworkMode(flags.NetworkName()),
		Init:         &trueValue,
	})
}

var trueValue = true
