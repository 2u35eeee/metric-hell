package content

import (
	"strings"
	"testing"

	"metric-hell/pkg/game"
)

func TestLoadNodesReadsDataFile(t *testing.T) {
	nodes, err := LoadNodes("../../data/nodes.json")
	if err != nil {
		t.Fatalf("LoadNodes returned error: %v", err)
	}
	if len(nodes) != 19 {
		t.Fatalf("len(nodes) = %d, want 19", len(nodes))
	}
	if nodes[0].ID != "gaokao_score" {
		t.Fatalf("first node ID = %q, want gaokao_score", nodes[0].ID)
	}
}

func TestNodesDescribeScenarioMeasurementAndBranches(t *testing.T) {
	nodes, err := LoadNodes("../../data/nodes.json")
	if err != nil {
		t.Fatalf("LoadNodes returned error: %v", err)
	}

	for _, node := range nodes {
		if node.Scenario == "" {
			t.Fatalf("node %s scenario is empty", node.ID)
		}
		if node.Measurement == "" {
			t.Fatalf("node %s measurement is empty", node.ID)
		}
		if len(node.Branches) < 2 {
			t.Fatalf("node %s has %d branches, want at least 2", node.ID, len(node.Branches))
		}
		for _, branch := range node.Branches {
			if branch.ID == "" || branch.Label == "" || branch.Scene == "" || branch.Description == "" {
				t.Fatalf("node %s has incomplete branch: %#v", node.ID, branch)
			}
		}
	}

	gpa := findNode(nodes, "gpa_loop")
	if gpa == nil {
		t.Fatal("gpa_loop node not found")
	}
	if !containsAll(gpa.Measurement, "4.0", "5.0") {
		t.Fatalf("gpa measurement = %q, want both 4.0 and 5.0 ranges", gpa.Measurement)
	}
}

func findNode(nodes []game.Node, id string) *game.Node {
	for _, node := range nodes {
		if node.ID == id {
			return &node
		}
	}
	return nil
}

func containsAll(value string, needles ...string) bool {
	for _, needle := range needles {
		if !strings.Contains(value, needle) {
			return false
		}
	}
	return true
}
