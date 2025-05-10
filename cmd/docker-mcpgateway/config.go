package main

import (
	"context"
	"encoding/json"
	"io"
	"net"
	"net/http"
	"os"
	"path/filepath"
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
	if err := get(ctx, httpClient(), "/volume-file-content?volumeId=docker-prompts&targetPath=registry.yaml", &content); err != nil {
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

func httpClient() *http.Client {
	return &http.Client{
		Transport: &http.Transport{
			DialContext: func(ctx context.Context, _, _ string) (conn net.Conn, err error) {
				return dialVolumeContents(ctx)
			},
		},
	}
}

func dialVolumeContents(ctx context.Context) (net.Conn, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	dialer := net.Dialer{}
	return dialer.DialContext(ctx, "unix", filepath.Join(home, "Library/Containers/com.docker.docker/Data/volume-contents.sock"))
}

func get(ctx context.Context, httpClient *http.Client, endpoint string, v any) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://localhost"+endpoint, nil)
	if err != nil {
		return err
	}
	req.Header.Set("X-DockerDesktop-Host", "vm.docker.internal")

	response, err := httpClient.Do(req)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	buf, err := io.ReadAll(response.Body)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(buf, &v); err != nil {
		return err
	}

	return nil
}
