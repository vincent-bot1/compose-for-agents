#!/bin/bash

# Test script for the Writer agent
# This script tests the Writer agent with a simple message

set -e  # Exit on any error

echo "Testing Writer agent..."
echo "Waiting for the service to be ready..."

# Wait for service to be ready (up to 60 seconds)
for i in {1..30}; do
  if curl -s "http://localhost:7777/v1/playground/status" > /dev/null; then
    echo "Service is ready!"
    break
  fi
  if [ $i -eq 30 ]; then
    echo "Timeout waiting for service to be ready"
    exit 1
  fi
  echo "Waiting... ($i/30)"
  sleep 2
done

# Get the Writer agent ID
echo "Getting Writer agent ID..."
WRITER_ID=$(curl -s 'http://localhost:7777/v1/playground/agents' | jq -r '.[] | select(.name=="Writer") | .agent_id')
if [ -z "$WRITER_ID" ]; then
  echo "Error: Writer agent not found"
  exit 1
fi
echo "Found Writer agent with ID: $WRITER_ID"

# Send a message to the agent and get the response
echo "Sending test message to Writer agent..."
RESPONSE=$(curl -s "http://localhost:7777/v1/playground/agents/$WRITER_ID/runs" -F "message=Hello" | jq 'select(.event=="RunCompleted")')

# Check if response is empty
if [ -z "$RESPONSE" ]; then
  echo "Error: Empty response from agent"
  exit 1
fi

# Extract and display the response content
CONTENT=$(echo "$RESPONSE" | jq -r '.content')

# Check if content is empty
if [ -z "$CONTENT" ] || [ "$CONTENT" = "null" ]; then
  echo "Error: Empty content in response"
  echo "Full response: $RESPONSE"
  exit 1
fi

# Check for error messages in content
if echo "$CONTENT" | grep -i "error" > /dev/null; then
  echo "Warning: Response contains error message"
  echo "Content: $CONTENT"
  exit 1
fi

# Extract and verify the model used
MODEL_USED=$(echo "$RESPONSE" | jq -r '.model')
echo "Model used: $MODEL_USED"

# Verify that a non-OpenAI model is being used
if [[ "$MODEL_USED" == *"gpt"* ]] || [[ "$MODEL_USED" == *"o1"* ]] || [[ "$MODEL_USED" == *"o3"* ]] || [[ "$MODEL_USED" == *"davinci"* ]] || [[ "$MODEL_USED" == *"curie"* ]] || [[ "$MODEL_USED" == *"babbage"* ]] || [[ "$MODEL_USED" == *"ada"* ]]; then
  echo "❌ FAIL: OpenAI model detected: $MODEL_USED"
  echo "Expected a non-OpenAI model (like qwen3, deepseek, gemma, etc.)"
  exit 1
else
  echo "✅ SUCCESS: Non-OpenAI model confirmed: $MODEL_USED"
fi

echo ""
echo "Agent response:"
echo "$CONTENT"
echo ""
echo "🎉 Test completed successfully!
