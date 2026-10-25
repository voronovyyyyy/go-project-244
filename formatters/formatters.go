package formatters

import (
	"fmt"
	"strconv"

	"code/internal/diff"
)

func Format(tree []diff.Node, name string) (string, error) {
	switch name {
	case "", "stylish":
		return stylish(tree, 0), nil
	case "plain":
		return plain(tree, ""), nil
	case "json":
		return formatJSON(tree)
	default:
		return "", fmt.Errorf("unsupported output format %q (expected stylish, plain or json)", name)
	}
}

func scalar(value any) string {
	switch typed := value.(type) {
	case nil:
		return "null"
	case float64:
		return strconv.FormatFloat(typed, 'f', -1, 64)
	default:
		return fmt.Sprint(value)
	}
}
