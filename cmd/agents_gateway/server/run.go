package server

import (
	"context"
	"fmt"
	"net"
	"os"
	"strings"
	"time"

	"github.com/docker/compose-agents-demo/pkg/catalog"
	"github.com/docker/compose-agents-demo/pkg/config"
	"github.com/docker/compose-agents-demo/pkg/docker"
	"github.com/docker/compose-agents-demo/pkg/sockets"
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

	// Read the MCP catalog.
	mcpCatalog := catalog.Get()

	// In standalone, ignore the registry.yaml passed on the command line
	// and read it from the docker volume.
	if g.Standalone {
		var err error
		g.RegistryYaml, err = docker.ReadPromptFile(ctx, "registry.yaml")
		if err != nil {
			return err
		}
	}

	// Which servers are enabled in the registry.yaml?
	registryConfig, err := config.ParseConfig(g.RegistryYaml)
	if err != nil {
		return fmt.Errorf("reading configuration: %w", err)
	}
	serverNames := registryConfig.ServerNames()

	// Detect which docker images are used.
	uniqueDockerImages := map[string]bool{}
	for _, serverName := range serverNames {
		serverSpec, tools, found := mcpCatalog.Find(serverName)

		switch {
		case !found:
			log("MCP server not found:", serverName)
		case serverSpec != nil:
			uniqueDockerImages[serverSpec.Image] = true
		case tools != nil:
			for _, tool := range *tools {
				uniqueDockerImages[tool.Container.Image] = true
			}
		}
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
	log("Pulling docker images", dockerImages)
	if err := client.PullImages(ctx, dockerImages...); err != nil {
		return fmt.Errorf("pulling docker images: %w", err)
	}
	log("Docker images pulled in", time.Since(startPull))

	// Then verify them. (TODO: should we check them, get the digest and pull that digest instead?)
	if g.VerifySignatures {
		if err := VerifySignatures(ctx, mcpImages); err != nil {
			return fmt.Errorf("verifying docker images: %w", err)
		}
	}

	// List all the available tools.
	startList := time.Now()
	log("Listing MCP tools...")
	serverTools, err := g.listTools(ctx, mcpCatalog, registryConfig, serverNames)
	if err != nil {
		return fmt.Errorf("listing tools: %w", err)
	}
	log(len(serverTools), "MCP tools listed in", time.Since(startList))

	toolCallbacks := callbacks(g.LogCalls, g.ScanSecrets)

	newStdioServer := func() *server.StdioServer {
		mcpServer := server.NewMCPServer("Docker AI MCP Gateway", "1.0.1", server.WithToolHandlerMiddleware(toolCallbacks))
		mcpServer.SetTools(serverTools...)
		return server.NewStdioServer(mcpServer)
	}

	log("Initialized MCP server in", time.Since(start))

	// Start the server
	if g.Standalone {
		return newStdioServer().Listen(ctx, os.Stdin, os.Stdout)
	}

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			conn, err := sockets.AcceptWithContext(ctx, ln)
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
