package main

import (
	"context"
	"flag"
	"log"
	"os/signal"
	"strings"
	"syscall"

	"github.com/docker/compose-agents-demo/cmd/agents_gateway/server"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	servers := flag.String("servers", "", "Comma-separated list of servers to enable")
	tools := flag.String("tools", "", "Comma-separated list of tools to enable")
	logCalls := flag.Bool("log_calls", false, "Log the tool calls")
	scanSecrets := flag.Bool("scan_secrets", false, "Verify that secrets are not passed to tools")
	verifyImages := flag.Bool("verify_images", false, "Verify the signatures off the images")
	flag.Parse()

	serverNames := parseCommaSeparated(*servers)
	toolNames := parseCommaSeparated(*tools)
	if err := server.Run(ctx, serverNames, toolNames, *logCalls, *scanSecrets, *verifyImages); err != nil {
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
