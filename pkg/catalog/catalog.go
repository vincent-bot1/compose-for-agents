package catalog

import (
	_ "embed"
	"fmt"

	"gopkg.in/yaml.v3"
)

//go:embed mcp-servers.yaml
var McpServersYAML []byte

//go:embed mcp-tools.yaml
var ToolsYAML []byte

func Get() (map[string]Server, []ToolGroup, error) {
	var servers []Server
	if err := yaml.Unmarshal(McpServersYAML, &servers); err != nil {
		return nil, nil, fmt.Errorf("reading servers catalog: %w", err)
	}

	byName := make(map[string]Server)
	for _, server := range servers {
		byName[server.Name] = server
	}

	var toolGroups []ToolGroup
	if err := yaml.Unmarshal(ToolsYAML, &toolGroups); err != nil {
		return nil, nil, fmt.Errorf("reading tools catalog: %w", err)
	}

	return byName, toolGroups, nil
}
