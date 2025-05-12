package config

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
	Config map[string]any `yaml:"config"`
}

// TODO(dga): I wanted to use the volume contents socket but in cloud mode, it isn't talking to the local Docker anymore.
func ReadPromptFile(ctx context.Context, name string) (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	path := filepath.Join(home, "Library/Containers/com.docker.docker/Data/docker.raw.sock")
	out, err := exec.CommandContext(ctx, "docker", "-H", "unix://"+path, "run", "--rm", "-v", "docker-prompts:/docker-prompts", "-w", "/docker-prompts", "busybox", "cat", name).Output()
	if err != nil {
		return "", err
	}

	return string(out), nil
}

func ParseConfig(registryYaml string) (Registry, error) {
	var registry Registry
	if err := yaml.Unmarshal([]byte(registryYaml), &registry); err != nil {
		return Registry{}, err
	}

	return registry, nil
}
