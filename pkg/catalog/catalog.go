package catalog

import (
	_ "embed"

	"gopkg.in/yaml.v3"
)

//go:embed mcp-servers.yaml
var McpServersYAML []byte

//go:embed mcp-tools.yaml
var ToolsYAML []byte

func Get() Catalog {
	// Guaranteed to parse, because we embed the yaml and have a test
	var servers []Server
	if err := yaml.Unmarshal(McpServersYAML, &servers); err != nil {
		panic(err)
	}

	serversByName := make(map[string]Server)
	for _, server := range servers {
		serversByName[server.Name] = server
	}

	// Guaranteed to parse, because we embed the yaml and have a test
	var toolGroups []ToolGroup
	if err := yaml.Unmarshal(ToolsYAML, &toolGroups); err != nil {
		panic(err)
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
	}
}
