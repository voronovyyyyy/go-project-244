package formatters

import (
	"fmt"
	"strings"

	"code/internal/diff"
)

func plain(tree []diff.Node, parent string) string {
	lines := make([]string, 0, len(tree))
	for _, node := range tree {
		path := node.Key
		if parent != "" {
			path = parent + "." + path
		}
		if line := plainNode(node, path); line != "" {
			lines = append(lines, line)
		}
	}
	return strings.Join(lines, "\n")
}

func plainNode(node diff.Node, path string) string {
	switch node.Type {
	case diff.Nested:
		return plain(node.Children, path)
	case diff.Added:
		return fmt.Sprintf("Property '%s' was added with value: %s", path, plainValue(node.After))
	case diff.Removed:
		return fmt.Sprintf("Property '%s' was removed", path)
	case diff.Updated:
		return fmt.Sprintf("Property '%s' was updated. From %s to %s", path, plainValue(node.Before), plainValue(node.After))
	default:
		return ""
	}
}

func plainValue(value any) string {
	switch typed := value.(type) {
	case map[string]any, []any:
		return "[complex value]"
	case string:
		return "'" + typed + "'"
	default:
		return scalar(value)
	}
}
