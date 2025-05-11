package catalog

import (
	"gopkg.in/yaml.v3"
)

type Catalog struct {
	Servers map[string]Server
	Tools   map[string]map[string]Tool
}

type Server struct {
	Name   string `yaml:"name" json:"name"`
	Image  string `yaml:"image" json:"image"`
	Run    Run    `yaml:"run,omitempty" json:"run,omitempty"`
	Config Config `yaml:"config,omitempty" json:"config,omitempty"`
}

type Secret struct {
	Id       string `yaml:"id" json:"id"`
	Name     string `yaml:"name" json:"name"`
	Value    string `yaml:"value" json:"value"`
	Required *bool  `yaml:"required,omitempty" json:"required,omitempty"`
}

type Env struct {
	Name       string `yaml:"name" json:"name"`
	Default    any    `yaml:"default" json:"default"`
	Expression string `yaml:"expression" json:"expression"`
}

type AnyOf struct {
	Required []string `yaml:"required,omitempty" json:"required,omitempty"`
}

type Schema struct {
	Type        string     `yaml:"type" json:"type"`
	Description string     `yaml:"description,omitempty" json:"description,omitempty"`
	Properties  SchemaList `yaml:"properties,omitempty" json:"properties,omitempty"`
	Required    []string   `yaml:"required,omitempty" json:"required,omitempty"`
	Items       Items      `yaml:"items,omitempty" json:"items,omitempty"`
	AnyOf       []AnyOf    `yaml:"anyOf,omitempty" json:"anyOf,omitempty"`
	Default     any        `yaml:"default,omitempty" json:"default,omitempty"`
}

type Items struct {
	Type string `yaml:"type" json:"type"`
}

type Run struct {
	Workdir string            `yaml:"workdir,omitempty" json:"workdir,omitempty"`
	Command []string          `yaml:"command,omitempty" json:"command,omitempty"`
	Volumes []string          `yaml:"volumes,omitempty" json:"volumes,omitempty"`
	Env     map[string]string `yaml:"env,omitempty" json:"env,omitempty"`
}

type Config struct {
	Description string   `yaml:"description,omitempty" json:"description,omitempty"`
	Secrets     []Secret `yaml:"secrets,omitempty" json:"secrets,omitempty"`
	Env         []Env    `yaml:"env,omitempty" json:"env,omitempty"`
	Parameters  Schema   `yaml:"parameters,omitempty" json:"parameters,omitempty"`
	AnyOf       []AnyOf  `yaml:"anyOf,omitempty" json:"anyOf,omitempty"`
}

type SchemaEntry struct {
	Schema Schema `yaml:",inline"`
	Name   string `yaml:"name"`
}

type SchemaList []SchemaEntry

func (tl *SchemaList) UnmarshalYAML(value *yaml.Node) error {
	for i := 0; i < len(value.Content); i += 2 {
		keyNode := value.Content[i]
		valNode := value.Content[i+1]

		var name string
		if err := keyNode.Decode(&name); err != nil {
			return err
		}

		var schema Schema
		if err := valNode.Decode(&schema); err != nil {
			return err
		}

		*tl = append(*tl, SchemaEntry{
			Name:   name,
			Schema: schema,
		})
	}
	return nil
}

type ToolGroup struct {
	Name  string `yaml:"name" json:"name"`
	Tools []Tool `yaml:"tools" json:"tools"`
}

type Tool struct {
	Name        string     `yaml:"name" json:"name"`
	Description string     `yaml:"description" json:"description"`
	Container   Container  `yaml:"container" json:"container"`
	Parameters  Parameters `yaml:"parameters" json:"parameters"`
}

type Parameters struct {
	Type       string     `yaml:"type" json:"type"`
	Properties Properties `yaml:"properties" json:"properties"`
}

type Properties map[string]Property

type Property struct {
	Type        string `yaml:"type" json:"type"`
	Description string `yaml:"description" json:"description"`
}

type Container struct {
	Image   string   `yaml:"image" json:"image"`
	Command []string `yaml:"command" json:"command"`
	Volumes []string `yaml:"volumes" json:"volumes"`
}
