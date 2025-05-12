package server

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/exec"
	"sort"
	"strings"
	"time"

	"github.com/docker/compose-agents-demo/pkg/catalog"
	"github.com/docker/compose-agents-demo/pkg/config"
	"github.com/docker/compose-agents-demo/pkg/docker"
	"github.com/mark3labs/mcp-go/server"
)

type Gateway struct {
	RegistryYaml     string
	ToolsNames       []string
	LogCalls         bool
	ScanSecrets      bool
	VerifySignatures bool
	Port             int
	Standalone       bool
}

func (g *Gateway) Run(ctx context.Context) error {
	start := time.Now()

	// Listen as early as possible to not lose client connections.
	var ln net.Listener
	if !g.Standalone {
		var (
			lc  net.ListenConfig
			err error
		)
		ln, err = lc.Listen(ctx, "tcp", fmt.Sprintf(":%d", g.Port))
		if err != nil {
			return err
		}
	}

	// Create docker client.
	client, err := docker.NewClient(ctx)
	if err != nil {
		return err
	}

	// In standalone, ignore the registry.yaml passed on the command line
	// and read it from the docker volume.
	if g.Standalone {
		var err error
		g.RegistryYaml, err = docker.ReadPromptFile(ctx, "registry.yaml")
		if err != nil {
			return err
		}
	}

	registryConfig, err := config.ParseConfig(g.RegistryYaml)
	if err != nil {
		return fmt.Errorf("reading configuration: %w", err)
	}

	// Read the MCP catalog.
	mcpCatalog, err := catalog.Get()
	if err != nil {
		return fmt.Errorf("listing catalog: %w", err)
	}

	// Which servers are enabled in the registry.yaml?
	var serverNames []string
	for serverName := range registryConfig.Servers {
		serverNames = append(serverNames, serverName)
	}
	sort.Strings(serverNames)

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

		fmt.Fprintln(os.Stderr, "MCP server not found:", serverName)
	}

	var (
		dockerImages []string
		mcpImages    []string
	)
	for image := range uniqueDockerImages {
		dockerImages = append(dockerImages, image)
		if strings.HasPrefix(image, "mcp/") {
			mcpImages = append(mcpImages, image)
		}
	}

	// Pull docker images first
	startPull := time.Now()
	fmt.Fprintln(os.Stderr, "Pulling docker images", dockerImages)
	if err := client.PullImages(ctx, dockerImages...); err != nil {
		return fmt.Errorf("pulling docker images: %w", err)
	}
	fmt.Fprintln(os.Stderr, "Docker images pulled in", time.Since(startPull))

	// Then verify them. (TODO: should we check them, get the digest and pull that digest instead?)
	if g.VerifySignatures {
		fmt.Fprintln(os.Stderr, "Verifying docker images", mcpImages)
		args := []string{"verify"}
		args = append(args, mcpImages...)
		args = append(args, "--key", "https://raw.githubusercontent.com/docker/keyring/refs/heads/main/public/mcp/latest.pub")

		cmd := exec.CommandContext(ctx, "/usr/bin/cosign", args...)
		cmd.Env = []string{"COSIGN_REPOSITORY=mcp/signatures"}
		if out, err := cmd.CombinedOutput(); err != nil {
			fmt.Fprintln(os.Stderr, "Failed to verify docker images:", string(out))
			return fmt.Errorf("verifying images: %w", err)
		}
		fmt.Fprintln(os.Stderr, "Docker images verified")
	}

	// List all the available tools.
	startList := time.Now()
	fmt.Fprintln(os.Stderr, "Listing MCP tools...")
	serverTools, err := g.listTools(ctx, mcpCatalog, registryConfig, serverNames)
	if err != nil {
		return fmt.Errorf("listing tools: %w", err)
	}
	fmt.Fprintln(os.Stderr, len(serverTools), "MCP tools listed in", time.Since(startList))

	toolCallbacks := callbacks(g.LogCalls, g.ScanSecrets)

	newStdioServer := func() *server.StdioServer {
		mcpServer := server.NewMCPServer("Docker AI MCP Gateway", "1.0.1", server.WithToolHandlerMiddleware(toolCallbacks))
		mcpServer.SetTools(serverTools...)
		return server.NewStdioServer(mcpServer)
	}

	fmt.Fprintln(os.Stderr, "Initialized MCP server in", time.Since(start))

	// Start the server
	if g.Standalone {
		return newStdioServer().Listen(ctx, os.Stdin, os.Stdout)
	}

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

				if err := newStdioServer().Listen(ctx, conn, conn); err != nil {
					fmt.Printf("Error listening: %v\n", err)
				}
			}()
		}
	}
}
