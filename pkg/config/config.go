package config

import (
	"context"
	"sort"

	"github.com/docker/compose-agents-demo/pkg/docker"
	"gopkg.in/yaml.v3"
)

type Registry struct {
	Servers map[string]Tile `yaml:"registry"`
}

func (r *Registry) ServerNames() []string {
	var names []string

	for name := range r.Servers {
		names = append(names, name)
	}
	sort.Strings(names)

	return names
}

type Tile struct {
	Config map[string]any `yaml:"config"`
}

func ReadPromptFile(ctx context.Context, name string) (string, error) {
	// Make sure to always talk to Docker Desktop directly in order to read the "local" volumes, those used by the MCP Toolkit extension.
	return docker.RunOnDockerDesktop(ctx, "-v", "docker-prompts:/docker-prompts", "-w", "/docker-prompts", "busybox", "cat", name)
}

func ParseConfig(registryYaml string) (Registry, error) {
	var registry Registry
	if err := yaml.Unmarshal([]byte(registryYaml), &registry); err != nil {
		return Registry{}, err
	}

	return registry, nil
}
