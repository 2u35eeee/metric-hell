package content

import (
	"encoding/json"
	"fmt"
	"os"

	"metric-hell/internal/game"
)

func LoadNodes(path string) ([]game.Node, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read nodes: %w", err)
	}

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
		seen[node.ID] = struct{}{}
	}
	return nil
}
