package server

import (
	"context"
	"fmt"
	"net"
	"os/exec"
	"strings"

	"github.com/docker/compose-agents-demo/pkg/catalog"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"golang.org/x/sync/errgroup"
)

func Run(ctx context.Context, serverNames, toolsNames []string, logCalls, scanSecrets, verifyImages bool) error {
	// Listen as early as possible to not lose client connections.
	var lc net.ListenConfig
	ln, err := lc.Listen(ctx, "tcp", ":8811")
	if err != nil {
		return err
	}

	// Read the MCP catalog.
	mcpCatalog, err := catalog.Get()
	if err != nil {
		return fmt.Errorf("listing catalog: %w", err)
	}

	// Detect which docker images are used.
	uniqueDockerImages := map[string]bool{}
	for _, serverName := range serverNames {
		// Is it an MCP Server?
		server, ok := mcpCatalog.Servers[serverName]
		if ok {
			uniqueDockerImages[server.Image] = true
			continue

		}

		// Is it a tool group?
		tools, ok := mcpCatalog.Tools[serverName]
		if ok {
			for _, tool := range tools {
				uniqueDockerImages[tool.Container.Image] = true
			}
			continue
		}

		fmt.Println("MCP server not found:", serverName)
	}
	var dockerImages []string
	for image := range uniqueDockerImages {
		dockerImages = append(dockerImages, image)
	}

	// Pull docker images first
	fmt.Println("Pulling docker images", dockerImages)
	var mcpImages []string
	errs, ctxPull := errgroup.WithContext(ctx)
	for _, dockerImage := range dockerImages {
		if strings.HasPrefix(dockerImage, "mcp/") {
			mcpImages = append(mcpImages, dockerImage)
		}

		errs.Go(func() error {
			cmd := exec.CommandContext(ctxPull, "docker", "pull", dockerImage)
			if err := cmd.Run(); err != nil {
				return fmt.Errorf("pulling docker image %s: %w", dockerImage, err)
			}
			return nil
		})
	}
	if err := errs.Wait(); err != nil {
		return fmt.Errorf("pulling docker images: %w", err)
	}
	fmt.Println("Docker images pulled")

	// Then verify them. (TODO: should we check them, get the digest and pull that digest instead?)
	if verifyImages {
		fmt.Println("Verifying docker images", mcpImages)
		args := []string{"verify"}
		args = append(args, mcpImages...)
		args = append(args, "--key", "https://raw.githubusercontent.com/docker/keyring/refs/heads/main/public/mcp/latest.pub")

		cmd := exec.CommandContext(ctx, "/usr/bin/cosign", args...)
		cmd.Env = []string{"COSIGN_REPOSITORY=mcp/signatures"}
		if out, err := cmd.CombinedOutput(); err != nil {
			fmt.Println("Failed to verify docker images:", string(out))
			return fmt.Errorf("verifying images: %w", err)
		}
		fmt.Println("Docker images verified")
	}

	// List all the available tools.
	serverTools, err := listTools(ctx, serverNames, mcpCatalog, toolsNames)
	if err != nil {
		return fmt.Errorf("listing tools: %w", err)
	}

	toolCallbacks := callbacks(logCalls, scanSecrets)

	// Server connections.
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			conn, err := acceptWithContext(ctx, ln)
			if err != nil {
				if ctx.Err() != nil {
					return ctx.Err()
				}
				fmt.Printf("Error accepting the connection: %v\n", err)
				continue
			}

			go func() {
				defer conn.Close()

				mcpServer := server.NewMCPServer("Docker AI MCP Gateway", "1.0.1", server.WithToolHandlerMiddleware(toolCallbacks))
				mcpServer.SetTools(serverTools...)
				stdioServer := server.NewStdioServer(mcpServer)

				if err := stdioServer.Listen(ctx, conn, conn); err != nil {
					fmt.Printf("Error listening: %v\n", err)
				}
			}()
		}
	}
}

func mcpServerHandler(server catalog.Server, tool mcp.Tool) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		client, err := startMCPClient(ctx, server, false)
		if err != nil {
			return nil, err
		}
		defer client.Close()

		return client.CallTool(ctx, tool.Name, request.Params.Arguments)
	}
}
