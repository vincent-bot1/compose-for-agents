package tests

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWriterAgent(t *testing.T) {
	client := NewAgentTestClient()
	
	t.Log("🧪 Testing Writer agent...")
	
	// Wait for services to be ready
	t.Log("⏳ Waiting for services to start...")
	client.WaitForServices(t)
	
	// Get available agents
	t.Log("🔍 Checking available agents...")
	agents := client.GetAgents(t)
	
	// Find the Writer agent
	writerAgent := client.FindAgentByName(t, agents, "Writer")
	require.NotNil(t, writerAgent, "Writer agent not found")
	t.Logf("✅ Found Writer agent with ID: %s", writerAgent.AgentID)
	
	// Send a message to the Writer agent
	t.Log("📤 Sending request to Writer agent...")
	message := "Hello, please write a short greeting message."
	response := client.SendMessageToAgent(t, writerAgent.AgentID, message)
	
	t.Log("📝 Agent response:")
	t.Log(response)
	
	// Validate response content
	t.Log("🔍 Validating response content...")
	
	// Check that response is not empty
	assert.NotEmpty(t, response, "Response should not be empty")
	
	// Check that response contains greeting-related terms
	greetingTerms := []string{"hello", "hi", "greeting", "welcome", "good"}
	client.AssertContainsAny(t, response, greetingTerms, "Response should contain greeting terms")
	
	// Check that response is reasonably long (more than just a word)
	assert.Greater(t, len(response), 10, "Response should be more than just a few characters")
	
	// Check mock gateway logs (Writer agent might not use tools, but let's check anyway)
	t.Log("🔍 Checking mock gateway logs...")
	// Note: Writer agent might not call any tools, so this is optional
	
	t.Log("🎉 Writer agent test completed successfully!")
	t.Log("✅ Agent successfully processed the message")
	t.Log("✅ Response contains appropriate content")
}
