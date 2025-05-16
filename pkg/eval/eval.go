package eval

import (
	"fmt"
	"reflect"
	"strings"
)

func Expression(expression string, config map[string]any) any {
	if !strings.HasPrefix(expression, "{{") || !strings.HasSuffix(expression, "}}") {
		return expression
	}

	path, rest, found := strings.Cut(expression[2:len(expression)-2], "|")
	value := dig(path, config)
	if !found {
		return value
	}

	for f := range strings.SplitSeq(rest, "|") {
		switch strings.TrimSpace(f) {
		case "volume":
			value = evaluate(value, volume)
		case "volume-target":
			value = evaluate(value, volumeTarget)
		case "safe", "into":
		}
	}

	return value
}

func evaluate(value any, fn func(v any) string) any {
	v := reflect.ValueOf(value)
	if v.Kind() == reflect.Slice {
		list := make([]string, v.Len())
		for i := range len(list) {
			list[i] = fn(v.Index(i).Interface())
		}
		return list
	}

	return fn(v)
}

func Expressions(expressions []string, arguments map[string]any) []string {
	var replaced []string

	for _, expression := range expressions {
		value := Expression(expression, arguments)

		v := reflect.ValueOf(value)
		if v.Kind() == reflect.Slice {
			for i := range v.Len() {
				replaced = append(replaced, fmt.Sprintf("%v", v.Index(i).Interface()))
			}
		} else {
			replaced = append(replaced, v.String())
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
