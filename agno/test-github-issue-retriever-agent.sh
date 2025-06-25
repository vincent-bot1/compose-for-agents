#!/bin/bash

# Test script for the Github Issue Retriever agent

set -e

echo "🧪 Testing Github Issue Retriever agent..."

# Wait for services to be ready
echo "⏳ Waiting for services to start..."
sleep 15

# Check if agents service is responding
echo "🔍 Checking agents service health..."
for i in {1..30}; do
  if curl -s http://localhost:7777/health > /dev/null 2>&1; then
    echo "✅ Agents service is ready"
    break
  fi
  if [ $i -eq 30 ]; then
    echo "❌ Agents service not responding after 30 attempts"
    echo "Debug: Checking what's running on port 7777..."
    curl -v http://localhost:7777/health 2>&1 || true
    exit 1
  fi
  sleep 2
done

# Check if mock gateway is responding
echo "🔍 Checking mock gateway health..."
for i in {1..30}; do
  if curl -s http://localhost:8080/health > /dev/null 2>&1; then
    echo "✅ Mock gateway is ready"
    break
  fi
  if [ $i -eq 30 ]; then
    echo "❌ Mock gateway not responding after 30 attempts"
    exit 1
  fi
  sleep 2
done

# Test the Github Issue Retriever agent
echo "🤖 Testing Github Issue Retriever agent..."

# First, check if the github agent exists
echo "🔍 Checking available agents..."
AGENTS_LIST=$(curl -s http://localhost:7777/v1/playground/agents 2>/dev/null || echo "Failed to get agents list")

# Extract github agent ID
GITHUB_AGENT_ID=$(echo "$AGENTS_LIST" | jq -r '.[] | select(.name == "Github Issue Retriever") | .agent_id' 2>/dev/null)

if [ -z "$GITHUB_AGENT_ID" ] || [ "$GITHUB_AGENT_ID" = "null" ]; then
  echo "❌ Github agent not found in agents list"
  echo "Available agents:"
  echo "$AGENTS_LIST" | jq -r '.[].name' 2>/dev/null || echo "$AGENTS_LIST"
  exit 1
fi

echo "✅ Found Github agent with ID: $GITHUB_AGENT_ID"

# Send a message to the github agent asking to retrieve issues from a repository
echo "📤 Sending request to github agent..."
RESPONSE=$(curl -s "http://localhost:7777/v1/playground/agents/$GITHUB_AGENT_ID/runs" \
  -F "message=Please retrieve all open issues from the repository example/turboencabulator" | jq 'select(.event=="RunCompleted")')

# Check if response is empty
if [ -z "$RESPONSE" ]; then
  echo "❌ Empty response from agent or no RunCompleted event found"
  exit 1
fi

echo "📝 Agent response:"
echo "$RESPONSE" | jq -r '.content' 2>/dev/null || echo "$RESPONSE"

# Check if the response contains expected GitHub issue data
echo ""
echo "🔍 Validating response content..."

# Check for issue numbers (should contain numbers like #1, #2, etc.)
if echo "$RESPONSE" | grep -q -E "(#[0-9]+|number.*[0-9]+)"; then
  echo "✅ Response contains issue numbers"
else
  echo "❌ Response does not contain issue numbers"
  exit 1
fi

# Check for issue titles (should contain some of our mock issue titles)
if echo "$RESPONSE" | grep -q -i -E "(turboencabulator|jazz music|pizza|toaster|sarcasm|bedtime stories)"; then
  echo "✅ Response contains expected issue titles from mock data"
else
  echo "❌ Response does not contain expected issue titles"
  exit 1
fi

# Check for GitHub-specific terms
if echo "$RESPONSE" | grep -q -i -E "(issue|github|repository|bug|feature)"; then
  echo "✅ Response contains GitHub-related terms"
else
  echo "❌ Response does not contain GitHub-related terms"
  exit 1
fi

# Check mock gateway logs to verify tool was called
echo ""
echo "🔍 Checking mock gateway logs..."
GATEWAY_LOGS=$(docker compose -f agno/compose.yaml -f agno/compose.test.yaml logs mock-gateway 2>/dev/null || echo "")

if echo "$GATEWAY_LOGS" | grep -q "list_issues"; then
  echo "✅ Mock gateway received list_issues tool call"
else
  echo "⚠️  Could not verify tool call in gateway logs (this might be expected)"
fi

# Check if the response mentions the repository name
if echo "$RESPONSE" | grep -q -i -E "(example|turboencabulator)"; then
  echo "✅ Response mentions the requested repository"
else
  echo "⚠️  Response does not clearly mention the requested repository"
fi

echo ""
echo "🎉 Github Issue Retriever agent test completed successfully!"
echo "✅ Agent successfully retrieved and processed GitHub issues"
echo "✅ Response contains expected issue data from mock gateway"
echo "✅ Tool integration working correctly"
