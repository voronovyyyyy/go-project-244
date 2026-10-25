package formatters

import (
	"sort"
	"strings"

	"code/internal/diff"
)

const indentWidth = 4

func stylish(tree []diff.Node, depth int) string {
	lines := []string{"{"}
	for _, node := range tree {
		lines = append(lines, stylishNode(node, depth+1)...)
	}
	lines = append(lines, indent(depth)+"}")
	return strings.Join(lines, "\n")
}

func stylishNode(node diff.Node, depth int) []string {
	line := func(sign, value string) string {
		return strings.Repeat(" ", depth*indentWidth-2) + sign + " " + node.Key + ": " + value
	}
	switch node.Type {
	case diff.Nested:
		return []string{line(" ", stylish(node.Children, depth))}
	case diff.Added:
		return []string{line("+", stylishValue(node.After, depth))}
	case diff.Removed:
		return []string{line("-", stylishValue(node.Before, depth))}
	case diff.Updated:
		return []string{line("-", stylishValue(node.Before, depth)), line("+", stylishValue(node.After, depth))}
	default:
		return []string{line(" ", stylishValue(node.Before, depth))}
	}
}

func stylishValue(value any, depth int) string {
	switch typed := value.(type) {
	case map[string]any:
		return stylishObject(typed, depth)
	case []any:
		items := make([]string, len(typed))
		for i, item := range typed {
			items[i] = stylishValue(item, depth)
		}
		return "[" + strings.Join(items, ", ") + "]"
	default:
		return scalar(value)
	}
}

func stylishObject(object map[string]any, depth int) string {
	if len(object) == 0 {
		return "{}"
	}
	keys := make([]string, 0, len(object))
	for key := range object {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	lines := []string{"{"}
	for _, key := range keys {
		lines = append(lines, indent(depth+1)+key+": "+stylishValue(object[key], depth+1))
	}
	lines = append(lines, indent(depth)+"}")
	return strings.Join(lines, "\n")
}

func indent(depth int) string {
	return strings.Repeat(" ", depth*indentWidth)
}
