package diff

import (
	"reflect"
	"sort"
)

type Type string

const (
	Added     Type = "added"
	Removed   Type = "deleted"
	Updated   Type = "changed"
	Unchanged Type = "unchanged"
	Nested    Type = "nested"
)

type Node struct {
	Key      string `json:"key"`
	Type     Type   `json:"type"`
	Before   any    `json:"value1,omitempty"`
	After    any    `json:"value2,omitempty"`
	Children []Node `json:"children,omitempty"`
}

func Build(first, second map[string]any) []Node {
	keys := make(map[string]struct{}, len(first)+len(second))
	for key := range first {
		keys[key] = struct{}{}
	}
	for key := range second {
		keys[key] = struct{}{}
	}
	sorted := make([]string, 0, len(keys))
	for key := range keys {
		sorted = append(sorted, key)
	}
	sort.Strings(sorted)
	result := make([]Node, 0, len(sorted))
	for _, key := range sorted {
		before, hasBefore := first[key]
		after, hasAfter := second[key]
		result = append(result, buildNode(key, before, after, hasBefore, hasAfter))
	}
	return result
}

func buildNode(key string, before, after any, hasBefore, hasAfter bool) Node {
	node := Node{Key: key, Before: before, After: after}
	oldObject, oldIsObject := before.(map[string]any)
	newObject, newIsObject := after.(map[string]any)
	switch {
	case !hasBefore:
		node.Type = Added
	case !hasAfter:
		node.Type = Removed
	case oldIsObject && newIsObject:
		node.Type = Nested
		node.Before, node.After = nil, nil
		node.Children = Build(oldObject, newObject)
	case reflect.DeepEqual(before, after):
		node.Type = Unchanged
		node.After = nil
	default:
		node.Type = Updated
	}
	return node
}
