package formatters

import (
	"encoding/json"

	"code/internal/diff"
)

type jsonDocument struct {
	Key      string      `json:"key"`
	Type     string      `json:"type"`
	Children []diff.Node `json:"children"`
}

func formatJSON(tree []diff.Node) (string, error) {
	result, err := json.MarshalIndent(jsonDocument{Key: "", Type: "root", Children: tree}, "", "  ")
	if err != nil {
		return "", err
	}
	return string(result), nil
}
