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

func Get() (Catalog, error) {
	var servers []Server
	if err := yaml.Unmarshal(McpServersYAML, &servers); err != nil {
		return Catalog{}, fmt.Errorf("reading servers catalog: %w", err)
	}

	serversByName := make(map[string]Server)
	for _, server := range servers {
		serversByName[server.Name] = server
	}

	var toolGroups []ToolGroup
	if err := yaml.Unmarshal(ToolsYAML, &toolGroups); err != nil {
		return Catalog{}, fmt.Errorf("reading tools catalog: %w", err)
	}

	toolsByName := make(map[string]map[string]Tool)
	for _, toolGroup := range toolGroups {
		toolsByName[toolGroup.Name] = map[string]Tool{}
		for _, tool := range toolGroup.Tools {
			toolsByName[toolGroup.Name][tool.Name] = tool
		}
	}

	return Catalog{
		Servers: serversByName,
		Tools:   toolsByName,
	}, nil
}
