package catalog

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGet(t *testing.T) {
	mcpCatalog := Get()

	assert.NotEmpty(t, mcpCatalog.Servers)
	assert.NotEmpty(t, mcpCatalog.Tools)
}

func TestServerImages(t *testing.T) {
	mcpCatalog := Get()

	for _, server := range mcpCatalog.Servers {
		assert.NotEmpty(t, server.Name)
		assert.NotEmpty(t, server.Image)
	}
	for name, tools := range mcpCatalog.Tools {
		assert.NotEmpty(t, name)
		for _, tool := range tools {
			assert.NotEmpty(t, tool.Name)
			assert.NotEmpty(t, tool.Container.Image)
		}
	}
}
