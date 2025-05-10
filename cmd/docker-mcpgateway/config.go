package main

import (
	"context"

	"gopkg.in/yaml.v3"
)

type FileContent struct {
	VolumeId   string `json:"volumeId"`
	TargetPath string `json:"targetPath"`
	Contents   string `json:"contents"`
}

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
	var content FileContent
	if err := get(ctx, httpClient(dialVolumeContents), "/volume-file-content?volumeId=docker-prompts&targetPath="+name, &content); err != nil {
		return "", err
	}

	return content.Contents, nil
}
