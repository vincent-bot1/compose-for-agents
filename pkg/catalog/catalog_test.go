package catalog

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGet(t *testing.T) {
	mcpCatalog, err := Get()

	assert.NotEmpty(t, mcpCatalog.Servers)
	assert.NotEmpty(t, mcpCatalog.Tools)
	assert.NoError(t, err)
}

func TestServerImages(t *testing.T) {
	mcpCatalog, err := Get()

	require.NoError(t, err)
	for _, server := range mcpCatalog.Servers {
		assert.NotEmpty(t, server.Name)
		assert.NotEmpty(t, server.Image)
	}

	require.NoError(t, err)
	for name, tools := range mcpCatalog.Tools {
		assert.NotEmpty(t, name)
		for _, tool := range tools {
			assert.NotEmpty(t, tool.Name)
			assert.NotEmpty(t, tool.Container.Image)
		}
	}
}
