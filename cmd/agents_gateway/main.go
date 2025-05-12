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

	registryYaml := flag.String("registry_yaml", "", "registry.yaml configuration file")
	tools := flag.String("tools", "", "Comma-separated list of tools to enable")
	logCalls := flag.Bool("log_calls", false, "Log the tool calls")
	scanSecrets := flag.Bool("scan_secrets", false, "Verify that secrets are not passed to tools")
	verifySignatures := flag.Bool("verify_signatures", false, "Verify the image signatures")
	port := flag.Int("port", 8811, "Port to listen on")
	standalone := flag.Bool("standalone", true, "Are we running in standalone mode?")

	// Parse flags and config
	flag.Parse()
	if *standalone && len(*registryYaml) > 0 {
		log.Fatalln("--registry_yaml is not supported in standalone mode")
	}

	gateway := server.Gateway{
		RegistryYaml:     *registryYaml,
		ToolsNames:       parseCommaSeparated(*tools),
		LogCalls:         *logCalls,
		ScanSecrets:      *scanSecrets,
		VerifySignatures: *verifySignatures,
		Port:             *port,
		Standalone:       *standalone,
	}
	if err := gateway.Run(ctx); err != nil {
		log.Fatalln(err)
	}
}

func parseCommaSeparated(values string) []string {
	var parsed []string

	for mcpImage := range strings.SplitSeq(values, ",") {
		name := strings.TrimSpace(mcpImage)
		if len(name) > 0 {
			parsed = append(parsed, strings.TrimSpace(mcpImage))
		}
	}

	return parsed
}
