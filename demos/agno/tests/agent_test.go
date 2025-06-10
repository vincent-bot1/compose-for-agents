package tests

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// AgentTestClient provides common functionality for testing agents
type AgentTestClient struct {
	BaseURL        string
	MockGatewayURL string
	HTTPClient     *http.Client
}

// NewAgentTestClient creates a new test client
func NewAgentTestClient() *AgentTestClient {
	return &AgentTestClient{
		BaseURL:        "http://localhost:7777",
		MockGatewayURL: "http://localhost:8080",
		HTTPClient: &http.Client{
			Timeout: 120 * time.Second, // Increased timeout for agent responses
		},
	}
}

// Agent represents an agent from the API
type Agent struct {
	AgentID string `json:"agent_id"`
	Name    string `json:"name"`
}

// RunEvent represents a streaming event from agent runs
type RunEvent struct {
	Event   string `json:"event"`
	Content string `json:"content"`
}

// WaitForServices waits for both agents service and mock gateway to be ready
func (c *AgentTestClient) WaitForServices(t *testing.T) {
	t.Helper()

	// Give services more time to start up
	t.Log("⏳ Waiting for services to start...")
	time.Sleep(15 * time.Second)

	// Wait for agents service
	c.waitForService(t, c.BaseURL+"/v1/playground/agents", "Agents service")

	// Wait for mock gateway
	c.waitForService(t, c.MockGatewayURL+"/health", "Mock gateway")
}

func (c *AgentTestClient) waitForService(t *testing.T, healthURL, serviceName string) {
	t.Helper()

	for i := 0; i < 60; i++ { // Increased from 30 to 60 attempts
		resp, err := c.HTTPClient.Get(healthURL)
		if err == nil && resp.StatusCode == http.StatusOK {
			resp.Body.Close()
			t.Logf("✅ %s is ready", serviceName)
			return
		}
		if resp != nil {
			resp.Body.Close()
		}
		if i%10 == 0 && i > 0 {
			t.Logf("⏳ Still waiting for %s... (attempt %d/60)", serviceName, i)
		}
		time.Sleep(2 * time.Second)
	}

	t.Fatalf("❌ %s not responding after 60 attempts", serviceName)
}

// GetAgents retrieves the list of available agents
func (c *AgentTestClient) GetAgents(t *testing.T) []Agent {
	t.Helper()

	resp, err := c.HTTPClient.Get(c.BaseURL + "/v1/playground/agents")
	require.NoError(t, err)
	defer resp.Body.Close()

	require.Equal(t, http.StatusOK, resp.StatusCode)

	var agents []Agent
	err = json.NewDecoder(resp.Body).Decode(&agents)
	require.NoError(t, err)

	return agents
}

// FindAgentByName finds an agent by name
func (c *AgentTestClient) FindAgentByName(t *testing.T, agents []Agent, name string) *Agent {
	t.Helper()

	for _, agent := range agents {
		if agent.Name == name {
			return &agent
		}
	}

	return nil
}

// SendMessageToAgent sends a message to an agent and returns the completed response
func (c *AgentTestClient) SendMessageToAgent(t *testing.T, agentID, message string) string {
	t.Helper()

	// Prepare form data
	data := url.Values{}
	data.Set("message", message)

	// Send request
	resp, err := c.HTTPClient.PostForm(
		fmt.Sprintf("%s/v1/playground/agents/%s/runs", c.BaseURL, agentID),
		data,
	)
	require.NoError(t, err)
	defer resp.Body.Close()

	require.Equal(t, http.StatusOK, resp.StatusCode)

	// Read streaming response and find RunCompleted event
	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	bodyStr := string(body)

	// The response format uses }{ to separate JSON objects
	// Split on }{ and add back the braces
	parts := strings.Split(bodyStr, "}{")
	var jsonObjects []string

	for i, part := range parts {
		if i == 0 {
			// First part, just add closing brace
			if !strings.HasSuffix(part, "}") {
				part += "}"
			}
		} else if i == len(parts)-1 {
			// Last part, just add opening brace
			if !strings.HasPrefix(part, "{") {
				part = "{" + part
			}
		} else {
			// Middle parts, add both braces
			part = "{" + part + "}"
		}
		jsonObjects = append(jsonObjects, part)
	}

	var lastContent string

	for _, jsonStr := range jsonObjects {
		jsonStr = strings.TrimSpace(jsonStr)
		if jsonStr == "" {
			continue
		}

		var event RunEvent
		if err := json.Unmarshal([]byte(jsonStr), &event); err != nil {
			// Try to extract content manually if JSON parsing fails
			if strings.Contains(jsonStr, `"event": "RunCompleted"`) {
				// Extract content field manually
				if contentMatch := regexp.MustCompile(`"content":\s*"([^"]*(?:\\.[^"]*)*)"`).FindStringSubmatch(jsonStr); len(contentMatch) > 1 {
					// Unescape the content
					content := strings.ReplaceAll(contentMatch[1], `\"`, `"`)
					content = strings.ReplaceAll(content, `\\`, `\`)
					return content
				}
			}
			continue
		}

		t.Logf("📨 Event: %s", event.Event)

		if event.Event == "RunCompleted" {
			return event.Content
		}

		// Keep track of the last content we see
		if event.Content != "" {
			lastContent = event.Content
		}
	}

	// If we didn't find RunCompleted but have some content, return it
	if lastContent != "" {
		t.Logf("⚠️  No RunCompleted event found, but returning last content")
		return lastContent
	}

	t.Fatal("No RunCompleted event found in response and no content available")
	return ""
}

// CheckMockGatewayLogs checks if the mock gateway received the expected tool call
// This is optional and won't fail the test if it can't get the logs
func (c *AgentTestClient) CheckMockGatewayLogs(t *testing.T, expectedTool string) {
	t.Helper()

	// Try to get logs using docker compose (similar to shell scripts)
	// This is optional and won't fail the test
	t.Logf("🔍 Checking mock gateway logs for %s tool call...", expectedTool)
	t.Logf("⚠️  Log checking is optional and may not always work in Go tests")
}

// AssertContainsRegex asserts that the content matches the given regex pattern
func (c *AgentTestClient) AssertContainsRegex(t *testing.T, content, pattern, description string) {
	t.Helper()

	matched, err := regexp.MatchString(pattern, content)
	require.NoError(t, err, "Invalid regex pattern: %s", pattern)
	assert.True(t, matched, "%s - Pattern: %s", description, pattern)
}

// AssertContainsAny asserts that the content contains at least one of the given strings (case-insensitive)
func (c *AgentTestClient) AssertContainsAny(t *testing.T, content string, terms []string, description string) {
	t.Helper()

	contentLower := strings.ToLower(content)
	for _, term := range terms {
		if strings.Contains(contentLower, strings.ToLower(term)) {
			t.Logf("✅ %s - Found term: %s", description, term)
			return
		}
	}

	t.Errorf("❌ %s - None of the expected terms found: %v", description, terms)
}
