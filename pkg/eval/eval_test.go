package eval

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEvaluateConstant(t *testing.T) {
	assert.Equal(t, "constant", Expression("constant", nil))
}

func TestEvaluate(t *testing.T) {
	assert.Equal(t, "value0", Expression("{{key0}}", map[string]any{"key0": "value0"}))
	assert.Equal(t, "value1", Expression("{{ key1 }}", map[string]any{"key1": "value1"}))
	assert.Equal(t, "value2", Expression("{{key2|safe}}", map[string]any{"key2": "value2"}))
}

func TestDotted(t *testing.T) {
	assert.Equal(t, "child_value0", Expression("{{top.key}}", map[string]any{"top": map[string]any{"key": "child_value0"}}))
	assert.Equal(t, "child_value1", Expression("{{top . key}}", map[string]any{"top": map[string]any{"key": "child_value1"}}))
	assert.Equal(t, "child_value2", Expression("{{top.key|ignored}}", map[string]any{"top": map[string]any{"key": "child_value2"}}))
}

func TestEvaluateUnknown(t *testing.T) {
	assert.Equal(t, "", Expression("{{unknown}}", nil))
	assert.Equal(t, "", Expression("{{top.unknown}}", map[string]any{"top": nil}))
}

func TestAtlassian(t *testing.T) {
	assert.Equal(t, "URL", Expression("{{atlassian.jira.url}}", map[string]any{"atlassian": map[string]any{"jira": map[string]any{"url": "URL", "username": "USERNAME"}}}))
}

func TestVolumes(t *testing.T) {
	assert.Equal(t, "/var/run/docker.sock:/var/run/docker.sock", Expression("/var/run/docker.sock:/var/run/docker.sock", nil))
	assert.Equal(t, []string{"/var/run/docker.sock:/var/run/docker.sock"}, Expressions([]string{"/var/run/docker.sock:/var/run/docker.sock"}, nil))
	assert.Equal(t, []string{"path1:path1", "path2:path2"}, Expressions([]string{"{{paths|volume|into}}"}, map[string]any{"paths": []string{"path1", "path2"}}))
	assert.Equal(t, []string{"path1", "path2"}, Expressions([]string{"{{paths|volume-targe|into}}"}, map[string]any{"paths": []string{"path1", "path2"}}))
}
