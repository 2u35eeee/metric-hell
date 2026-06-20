package content

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"os"

	"metric-hell/pkg/game"
)

func LoadNodes(path string) ([]game.Node, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read nodes: %w", err)
	}
	return parseNodes(data)
}

func LoadNodesFS(fsys fs.FS, path string) ([]game.Node, error) {
	data, err := fs.ReadFile(fsys, path)
	if err != nil {
		return nil, fmt.Errorf("read nodes: %w", err)
	}
	return parseNodes(data)
}

func parseNodes(data []byte) ([]game.Node, error) {
	var nodes []game.Node
	if err := json.Unmarshal(data, &nodes); err != nil {
		return nil, fmt.Errorf("parse nodes: %w", err)
	}
	if err := validateNodes(nodes); err != nil {
		return nil, err
	}
	return nodes, nil
}

func validateNodes(nodes []game.Node) error {
	if len(nodes) == 0 {
		return fmt.Errorf("nodes cannot be empty")
	}
	seen := make(map[string]struct{}, len(nodes))
	for i, node := range nodes {
		if node.ID == "" {
			return fmt.Errorf("node %d has empty id", i)
		}
		if node.Title == "" {
			return fmt.Errorf("node %s has empty title", node.ID)
		}
		if _, ok := seen[node.ID]; ok {
			return fmt.Errorf("duplicate node id %q", node.ID)
		}
		if node.Input.Type != game.InputTypeNumber && node.Input.Type != game.InputTypeSelect {
			return fmt.Errorf("node %s has invalid input type %q", node.ID, node.Input.Type)
		}
		if node.Input.Prompt == "" {
			return fmt.Errorf("node %s has empty input prompt", node.ID)
		}
		if len(node.Options) == 0 {
			return fmt.Errorf("node %s has no answer options", node.ID)
		}
		for j, option := range node.Options {
			if option.ID == "" {
				return fmt.Errorf("node %s option %d has empty id", node.ID, j)
			}
			if option.Label == "" {
				return fmt.Errorf("node %s option %s has empty label", node.ID, option.ID)
			}
			if option.Verdict == "" {
				return fmt.Errorf("node %s option %s has empty verdict", node.ID, option.ID)
			}
			if option.Proof == "" {
				return fmt.Errorf("node %s option %s has empty proof", node.ID, option.ID)
			}
		}
		seen[node.ID] = struct{}{}
	}
	return nil
}
