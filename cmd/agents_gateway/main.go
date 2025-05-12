package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os/signal"
	"strings"
	"syscall"

	"github.com/docker/compose-agents-demo/cmd/agents_gateway/server"
	"github.com/docker/compose-agents-demo/pkg/config"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	registryYaml := flag.String("registry_yaml", "", "registry.yaml configuration file")
	tools := flag.String("tools", "", "Comma-separated list of tools to enable")
	logCalls := flag.Bool("log_calls", false, "Log the tool calls")
	scanSecrets := flag.Bool("scan_secrets", false, "Verify that secrets are not passed to tools")
	verifySignatures := flag.Bool("verify_signatures", false, "Verify the image signatures")
	flag.Parse()

	// Parse flags and config
	registryConfig, err := config.ParseConfig(*registryYaml)
	if err != nil {
		log.Fatalln(fmt.Errorf("reading configuration: %w", err))
	}
	toolNames := parseCommaSeparated(*tools)

	if err := server.Run(ctx, registryConfig, toolNames, *logCalls, *scanSecrets, *verifySignatures); err != nil {
		log.Fatalln(err)
	}
}

func parseCommaSeparated(values string) []string {
	var parsed []string

	for mcpImage := range strings.SplitSeq(values, ",") {
		parsed = append(parsed, strings.TrimSpace(mcpImage))
	}

	return parsed
}
