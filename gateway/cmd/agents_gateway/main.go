package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/docker/gateway/cmd/agents_gateway/secrets"
	"github.com/mark3labs/mcp-go/client"
	"github.com/mark3labs/mcp-go/client/transport"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

const defaultMCPGatewayHost = "host.docker.internal:8811"

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	tools := flag.String("tools", "", "Comma-separated list of tools to enable")
	logCalls := flag.Bool("logCalls", false, "Log the tool calls")
	scanSecrets := flag.Bool("scanSecrets", false, "Verify that secrets are not passed to tools")
	flag.Parse()

	if err := run(ctx, *tools, *logCalls, *scanSecrets); err != nil {
		log.Fatalln(err)
	}
}

func run(ctx context.Context, tools string, logCalls, scanSecrets bool) error {
	toolNeeded := map[string]bool{}
	for tool := range strings.SplitSeq(tools, ",") {
		toolNeeded[tool] = true
	}

	c, err := startClient(ctx)
	if err != nil {
		return fmt.Errorf("starting client: %w", err)
	}
	defer c.Close()

	srv := &http.Server{
		Addr: ":8811",
		Handler: server.NewSSEServer(server.NewMCPServer(
			"Docker AI MCP Gateway",
			"1.0.1",
			server.WithToolCapabilities(true),
			server.WithListToolsHandler(func(ctx context.Context, request mcp.ListToolsRequest) (*mcp.ListToolsResult, error) {
				list, err := c.ListTools(ctx, request)
				if err != nil {
					return nil, err
				}

				var filtered []mcp.Tool
				for _, tool := range list.Tools {
					if len(toolNeeded) == 0 || toolNeeded[tool.Name] {
						filtered = append(filtered, tool)
					}
				}

				return &mcp.ListToolsResult{
					Tools: filtered,
				}, nil
			}),
			server.WithToolCallHandler(func(_ context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
				toolName := request.Params.Name
				if _, ok := toolNeeded[toolName]; !ok {
					return nil, fmt.Errorf("tool %s is not available", toolName)
				}

				// Print arguments into a string
				var arguments string
				buf, err := json.Marshal(request.Params.Arguments)
				if err != nil {
					arguments = fmt.Sprintf("%v", request.Params.Arguments)
				} else {
					arguments = string(buf)
				}

				// Callbacks
				if scanSecrets {
					fmt.Printf("Scanning tool call arguments for secrets...\n")

					if secrets.ContainsSecrets(arguments) {
						return nil, fmt.Errorf("a secret is being passed to tool %s", toolName)
					}
				}
				if logCalls {
					fmt.Printf("Calling tool %s with arguments: %s\n", toolName, arguments)
				}

				// Actual call
				return c.CallTool(ctx, request)
			}),
		)),
	}

	go func() {
		<-ctx.Done()
		fmt.Println("Shutting down")
		srv.Shutdown(context.Background())
	}()

	var lc net.ListenConfig
	ln, err := lc.Listen(ctx, "tcp", srv.Addr)
	if err != nil {
		return err
	}

	fmt.Println("Starting Agents Gateway on", srv.Addr)
	return srv.Serve(ln)
}

func startClient(ctx context.Context) (*client.Client, error) {
	host := os.Getenv("MCPGATEWAY_ENDPOINT")
	if host == "" {
		host = defaultMCPGatewayHost
	}

	conn, err := net.Dial("tcp", host)
	if err != nil {
		return nil, fmt.Errorf("dialing: %w", err)
	}

	c := client.NewClient(transport.NewIO(conn, conn, conn))
	if err := c.Start(ctx); err != nil {
		return nil, fmt.Errorf("starting client: %w", err)
	}

	initRequest := mcp.InitializeRequest{}
	initRequest.Params.ProtocolVersion = mcp.LATEST_PROTOCOL_VERSION
	initRequest.Params.ClientInfo = mcp.Implementation{
		Name:    "docker",
		Version: "1.0.0",
	}

	if _, err := c.Initialize(ctx, initRequest); err != nil {
		return nil, fmt.Errorf("initializing: %w", err)
	}

	return c, nil
}
