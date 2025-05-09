package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/mark3labs/mcp-go/client"
	"github.com/mark3labs/mcp-go/client/transport"
	"github.com/mark3labs/mcp-go/mcp"
)

func main() {
	args := os.Args[1:]
	if len(args) == 0 {
		usage()
		return
	}

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	switch args[0] {
	case "list":
		verbose := true
		if err := list(ctx, verbose); err != nil {
			log.Fatal(err)
		}
	case "count":
		verbose := false
		if err := list(ctx, verbose); err != nil {
			log.Fatal(err)
		}
	case "call":
		if err := call(ctx, args[1:]); err != nil {
			log.Fatal(err)
		}
	default:
		usage()
	}
}

func usage() {
	fmt.Println("Usage: client COMMAND [ARGS]")
	fmt.Println()
	fmt.Println("A command-line debug client for the MCP.")
	fmt.Println()
	fmt.Println("Commands:")
	fmt.Println("  list     List all available tools")
	fmt.Println("  count    Count all available tools")
	fmt.Println("  call     Call a specific tool with arguments")
}

func list(ctx context.Context, verbose bool) error {
	c, err := start(ctx)
	if err != nil {
		return fmt.Errorf("starting client: %w", err)
	}
	defer c.Close()

	response, err := c.ListTools(ctx, mcp.ListToolsRequest{})
	if err != nil {
		return fmt.Errorf("listing tools: %w", err)
	}

	buf, err := json.MarshalIndent(response.Tools, "", "  ")
	if err != nil {
		return fmt.Errorf("marshalling tools: %w", err)
	}

	if verbose {
		fmt.Println(len(response.Tools), "tools:")
		fmt.Println(string(buf))
	} else {
		fmt.Println(len(response.Tools), "tools")
	}

	return nil
}

func call(ctx context.Context, args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("no tool name provided")
	}
	toolName := args[0]

	c, err := start(ctx)
	if err != nil {
		return fmt.Errorf("starting client: %w", err)
	}
	defer c.Close()

	request := mcp.CallToolRequest{}
	request.Params.Name = toolName
	request.Params.Arguments = parseArgs(args[1:])

	start := time.Now()
	response, err := c.CallTool(ctx, request)
	if err != nil {
		return fmt.Errorf("calling tool: %w", err)
	}
	c.Close()
	fmt.Println("Tool call took:", time.Since(start))

	if response.IsError {
		return fmt.Errorf("error calling tool: %s", toolName)
	}

	for _, content := range response.Content {
		if textContent, ok := content.(mcp.TextContent); ok {
			fmt.Println(textContent.Text)
		} else {
			fmt.Println(content)
		}
	}

	return nil
}

func start(ctx context.Context) (*client.Client, error) {
	host := os.Getenv("MCPGATEWAY_ENDPOINT")
	conn, err := net.Dial("tcp", host)
	if err != nil {
		return nil, fmt.Errorf("dialing: %w", err)
	}

	c := client.NewClient(transport.NewIO(conn, conn, conn))
	if err := c.Start(ctx); err != nil {
		return nil, err
	}

	initRequest := mcp.InitializeRequest{}
	initRequest.Params.ProtocolVersion = mcp.LATEST_PROTOCOL_VERSION
	initRequest.Params.ClientInfo = mcp.Implementation{
		Name:    "docker",
		Version: "1.0.0",
	}

	ctx, cancel := context.WithTimeout(ctx, 20*time.Second)
	defer cancel()

	if _, err := c.Initialize(ctx, initRequest); err != nil {
		return nil, fmt.Errorf("initializing: %w", err)
	}

	return c, nil
}

func parseArgs(args []string) map[string]any {
	parsed := map[string]any{}

	for _, arg := range args {
		parts := strings.SplitN(arg, "=", 2)
		if len(parts) == 2 {
			parsed[parts[0]] = parts[1]
		} else {
			parsed[arg] = nil
		}
	}

	// MCP servers return an error if the args are empty so we make sure
	// there is at least one argument
	if len(parsed) == 0 {
		parsed["args"] = "..."
	}

	return parsed
}
