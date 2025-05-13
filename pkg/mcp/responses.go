package mcp

import (
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
)

func ToolError(errorText string) *mcp.CallToolResult {
	return &mcp.CallToolResult{
		Result: mcp.Result{
			Meta: map[string]any{},
		},
		Content: []mcp.Content{
			mcp.TextContent{
				Type: "text",
				Text: errorText,
			},
		},
		IsError: true,
	}
}

func ToolResult(text string) *mcp.CallToolResult {
	out := text
	if len(strings.TrimSpace(text)) == 0 {
		out = "There was no output from the tool call"
	}

	return &mcp.CallToolResult{
		Result: mcp.Result{
			Meta: map[string]any{},
		},
		Content: []mcp.Content{
			mcp.TextContent{
				Type: "text",
				Text: out,
			},
		},
	}
}
