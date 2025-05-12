package main

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type Registry struct {
	Servers map[string]Tile `yaml:"registry"`
}

type Tile struct {
	Config Config `yaml:"config"`
}

type Config map[string]map[string]any

func enabledMCPServers(ctx context.Context) (map[string]Tile, error) {
	content, err := readPromptFile(ctx, "registry.yaml")
	if err != nil {
		return nil, err
	}

	var registry Registry
	if err := yaml.Unmarshal([]byte(content), &registry); err != nil {
		return nil, err
	}

	return registry.Servers, nil
}

func readPromptFile(ctx context.Context, name string) (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	// TODO(dga): I wanted to use the volume contents socket but in cloud mode, it isn't talking to the local Docker anymore.
	path := filepath.Join(home, "Library/Containers/com.docker.docker/Data/docker.raw.sock")
	out, err := exec.CommandContext(ctx, "docker", "-H", "unix://"+path, "run", "--rm", "-v", "docker-prompts:/docker-prompts", "-w", "/docker-prompts", "busybox", "cat", name).Output()
	if err != nil {
		return "", err
	}

	return string(out), nil
}
