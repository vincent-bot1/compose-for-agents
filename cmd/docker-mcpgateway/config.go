package main

import (
	"context"
	"sort"

	"gopkg.in/yaml.v3"
)

type FileContent struct {
	VolumeId   string `json:"volumeId"`
	TargetPath string `json:"targetPath"`
	Contents   string `json:"contents"`
}

type Registry struct {
	Servers map[string]any `json:"registry" yaml:"registry"`
}

func enabledMCPServers(ctx context.Context) ([]string, error) {
	var content FileContent
	if err := get(ctx, httpClient(dialVolumeContents), "/volume-file-content?volumeId=docker-prompts&targetPath=registry.yaml", &content); err != nil {
		return nil, err
	}

	var registry Registry
	if err := yaml.Unmarshal([]byte(content.Contents), &registry); err != nil {
		return nil, err
	}

	var servers []string
	for server := range registry.Servers {
		servers = append(servers, server)
	}
	sort.Strings(servers)

	return servers, nil
}
