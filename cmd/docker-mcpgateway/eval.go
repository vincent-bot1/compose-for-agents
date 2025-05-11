package main

import (
	"fmt"
	"strings"
)

func evaluate(expression string, config map[string]any) string {
	if !strings.HasPrefix(expression, "{{") || !strings.HasSuffix(expression, "}}") {
		return expression
	}

	parts := strings.Split(expression[2:len(expression)-2], "|")

	return fmt.Sprintf("%s", dig(parts[0], config))
}

func dig(key string, config map[string]any) any {
	key = strings.TrimSpace(key)

	top, rest, found := strings.Cut(key, ".")
	if !found {
		value := config[key]
		if value == nil {
			return ""
		}
		return config[key]
	}

	top = strings.TrimSpace(top)
	childConfig, ok := config[top].(map[string]any)
	if !ok {
		return ""
	}

	return dig(rest, childConfig)
}
