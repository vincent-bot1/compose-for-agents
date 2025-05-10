package catalog

import (
	_ "embed"
	"fmt"

	"gopkg.in/yaml.v3"
)

//go:embed mcp-servers.yaml
var McpServersYAML []byte

func Get() (map[string]Server, error) {
	var servers []Server
	if err := yaml.Unmarshal(McpServersYAML, &servers); err != nil {
		return nil, fmt.Errorf("reading catalog: %w", err)
	}

	byName := make(map[string]Server)
	for _, server := range servers {
		byName[server.Name] = server
	}
	return byName, nil
}
