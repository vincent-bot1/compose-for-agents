package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"strings"
	"syscall"

	mcpclient "github.com/docker/gateway/cmd/agents_gateway/mcp"
	"github.com/docker/gateway/cmd/agents_gateway/secrets"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	servers := flag.String("servers", "", "Comma-separated list of servers to enable")
	config := flag.String("config", "", "Comma-separated list of config for the servers")
	tools := flag.String("tools", "", "Comma-separated list of tools to enable")
	logCalls := flag.Bool("log_calls", false, "Log the tool calls")
	scanSecrets := flag.Bool("scan_secrets", false, "Verify that secrets are not passed to tools")
	flag.Parse()

	if err := run(ctx, *servers, *config, *tools, *logCalls, *scanSecrets); err != nil {
		log.Fatalln(err)
	}
}

func run(ctx context.Context, servers, config, tools string, logCalls, scanSecrets bool) error {
	// List as early as possible to not lose client connections
	var lc net.ListenConfig
	ln, err := lc.Listen(ctx, "tcp", ":8811")
	if err != nil {
		return err
	}

	// Filter out tools
	toolNeeded := map[string]bool{}
	for tool := range strings.SplitSeq(tools, ",") {
		toolNeeded[strings.TrimSpace(tool)] = true
	}

	mcpServer := server.NewMCPServer("Docker AI MCP Gateway", "1.0.1", server.WithToolHandlerMiddleware(func(next server.ToolHandlerFunc) server.ToolHandlerFunc {
		return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			// Print arguments into a string
			var arguments string
			buf, err := json.Marshal(request.Params.Arguments)
			if err != nil {
				arguments = fmt.Sprintf("%v", request.Params.Arguments)
			} else {
				arguments = string(buf)
			}

			if scanSecrets {
				fmt.Printf("Scanning tool call arguments for secrets...\n")
				if secrets.ContainsSecrets(arguments) {
					return nil, fmt.Errorf("a secret is being passed to tool %s", request.Params.Name)
				}
			}

			if logCalls {
				fmt.Printf("Calling tool %s with arguments: %s\n", request.Params.Name, arguments)
			}

			return next(ctx, request)
		}
	}))

	for mcpImage := range strings.SplitSeq(servers, ",") {
		mcpImage := strings.TrimSpace(mcpImage)

		pull := true
		client, err := startMCPClient(ctx, mcpImage, pull, config)
		if err != nil {
			return err
		}

		tools, err := client.ListTools(ctx)
		client.Close()
		if err != nil {
			return fmt.Errorf("listing tools: %w", err)
		}

		for _, tool := range tools {
			if _, ok := toolNeeded[tool.Name]; !ok {
				continue
			}

			mcpServer.AddTool(tool, mcpServerHandler(mcpImage, tool, config))
		}
	}

	stdioServer := server.NewStdioServer(mcpServer)
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			conn, err := ln.Accept()
			if err != nil {
				fmt.Printf("Error accepting the connection: %v\n", err)
				continue
			}

			go func() {
				defer conn.Close()
				if err := stdioServer.Listen(ctx, conn, conn); err != nil {
					fmt.Printf("Error listening: %v\n", err)
				}
			}()
		}
	}
}

func mcpServerHandler(mcpImage string, tool mcp.Tool, config string) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		client, err := startMCPClient(ctx, mcpImage, false, config)
		if err != nil {
			return nil, err
		}
		defer client.Close()

		return client.CallTool(ctx, tool.Name, request.Params.Arguments)
	}
}

// config: mcp/github-mcp-server.GITHUB_PERSONAL_ACCESS_TOKEN=$GITHUB_TOKEN
func startMCPClient(ctx context.Context, mcpImage string, pull bool, config string) (*mcpclient.Client, error) {
	args := []string{}
	for cfg := range strings.SplitSeq(config, ",") {
		prefix := mcpImage + "."
		if !strings.HasPrefix(cfg, prefix) {
			continue
		}

		mapping := strings.TrimPrefix(cfg, prefix)
		parts := strings.SplitN(mapping, "=", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid config format: %s", cfg)
		}

		if strings.HasPrefix(parts[1], "$") {
			// TODO: find a better way to pass this secret
			args = append(args, "-e", parts[0]+"="+os.Getenv(parts[1][1:]))
		} else {
			args = append(args, "-e", parts[0]+"="+parts[1])
		}
	}

	client := mcpclient.NewClientArgs(mcpImage, pull, args, nil)
	if err := client.Start(ctx); err != nil {
		return nil, fmt.Errorf("failed to start server %s: %w", mcpImage, err)
	}

	return client, nil
}
