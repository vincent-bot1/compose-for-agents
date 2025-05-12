package eval

import (
	"fmt"
	"strings"
)

func Expression(expression string, config map[string]any) any {
	if !strings.HasPrefix(expression, "{{") || !strings.HasSuffix(expression, "}}") {
		return expression
	}

	parts := strings.Split(expression[2:len(expression)-2], "|")

	return fmt.Sprintf("%s", dig(parts[0], config))
}

func Expressions(expressions []string, arguments map[string]any) []string {
	var replaced []string

	for _, expression := range expressions {
		value := Expression(expression, arguments)

		switch v := value.(type) {
		case []any:
			for _, vv := range v {
				replaced = append(replaced, fmt.Sprintf("%v", vv))
			}
		default:
			replaced = append(replaced, fmt.Sprintf("%v", v))
		}
	}

	return replaced
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
