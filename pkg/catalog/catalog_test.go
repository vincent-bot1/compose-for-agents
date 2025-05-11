package catalog

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGet(t *testing.T) {
	servers, tools, err := Get()

	assert.NotEmpty(t, servers)
	assert.NotEmpty(t, tools)
	assert.NoError(t, err)
}

func TestServerImages(t *testing.T) {
	servers, _, err := Get()

	require.NoError(t, err)
	for _, server := range servers {
		assert.NotEmpty(t, server.Name)
		assert.NotEmpty(t, server.Image)
	}
}

func TestToolImages(t *testing.T) {
	_, groups, err := Get()

	require.NoError(t, err)
	for _, group := range groups {
		assert.NotEmpty(t, group.Name)
		for _, tool := range group.Tools {
			assert.NotEmpty(t, tool.Name)
			assert.NotEmpty(t, tool.Container.Image)
		}
	}
}
