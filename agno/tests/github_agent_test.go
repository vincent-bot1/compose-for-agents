package tests

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGitHubIssueRetrieverAgent(t *testing.T) {
	client := NewAgentTestClient()
	
	t.Log("🧪 Testing GitHub Issue Retriever agent...")
	
	// Wait for services to be ready
	t.Log("⏳ Waiting for services to start...")
	client.WaitForServices(t)
	
	// Get available agents
	t.Log("🔍 Checking available agents...")
	agents := client.GetAgents(t)
	
	// Find the GitHub Issue Retriever agent
	githubAgent := client.FindAgentByName(t, agents, "Github Issue Retriever")
	require.NotNil(t, githubAgent, "GitHub Issue Retriever agent not found")
	t.Logf("✅ Found GitHub agent with ID: %s", githubAgent.AgentID)
	
	// Send a message to the GitHub agent
	t.Log("📤 Sending request to GitHub agent...")
	message := "Please retrieve all open issues from the repository example/turboencabulator"
	response := client.SendMessageToAgent(t, githubAgent.AgentID, message)
	
	t.Log("📝 Agent response:")
	t.Log(response)
	
	// Validate response content
	t.Log("🔍 Validating response content...")
	
	// Check that response is not empty
	assert.NotEmpty(t, response, "Response should not be empty")
	
	// Check for issue numbers (should contain numbers like #16, #17, etc.)
	client.AssertContainsRegex(t, response, `#[0-9]+`, "Response should contain issue numbers")
	
	// Check for expected issue titles from mock data
	mockIssueTitles := []string{
		"turboencabulator", "jazz music", "pizza", "toaster", 
		"sarcasm", "bedtime stories", "pirate speak", "emojis",
	}
	client.AssertContainsAny(t, response, mockIssueTitles, "Response should contain expected issue titles from mock data")
	
	// Check for GitHub-specific terms
	githubTerms := []string{"issue", "github", "repository", "bug", "feature"}
	client.AssertContainsAny(t, response, githubTerms, "Response should contain GitHub-related terms")
	
	// Check if the response mentions the requested repository
	repoTerms := []string{"example", "turboencabulator"}
	client.AssertContainsAny(t, response, repoTerms, "Response should mention the requested repository")
	
	// Check mock gateway logs for list_issues tool call
	t.Log("🔍 Checking mock gateway logs...")
	client.CheckMockGatewayLogs(t, "list_issues")
	
	t.Log("🎉 GitHub Issue Retriever agent test completed successfully!")
	t.Log("✅ Agent successfully retrieved and processed GitHub issues")
	t.Log("✅ Response contains expected issue data from mock gateway")
	t.Log("✅ Tool integration working correctly")
}
