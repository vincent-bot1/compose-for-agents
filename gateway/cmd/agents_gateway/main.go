package main

import (
	"context"
	"flag"
	"log"
	"os/signal"
	"syscall"

	"github.com/docker/compose-agents-demo/gateway/cmd/agents_gateway/server"
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

	if err := server.Run(ctx, *servers, *config, *tools, *logCalls, *scanSecrets); err != nil {
		log.Fatalln(err)
	}
}
