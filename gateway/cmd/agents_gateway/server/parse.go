package server

import "strings"

func parseServers(servers string) []string {
	return parseCommaSeparated(servers)
}

func parseConfig(config string) []string {
	return parseCommaSeparated(config)
}

func parseCommaSeparated(values string) []string {
	var parsed []string

	for mcpImage := range strings.SplitSeq(values, ",") {
		parsed = append(parsed, strings.TrimSpace(mcpImage))
	}

	return parsed
}
