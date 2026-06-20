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

func TestNodesDescribeInputsOptionsAndAuditCopy(t *testing.T) {
	nodes, err := LoadNodes("../../data/nodes.json")
	if err != nil {
		t.Fatalf("LoadNodes returned error: %v", err)
	}

	for _, node := range nodes {
		if node.Input.Type != game.InputTypeNumber && node.Input.Type != game.InputTypeSelect {
			t.Fatalf("node %s input type = %q, want number or select", node.ID, node.Input.Type)
		}
		if node.Input.Prompt == "" {
			t.Fatalf("node %s input prompt is empty", node.ID)
		}
		if len(node.Options) < 2 {
			t.Fatalf("node %s has %d options, want at least 2", node.ID, len(node.Options))
		}
		for _, option := range node.Options {
			if option.ID == "" || option.Label == "" || option.Verdict == "" || option.Proof == "" {
				t.Fatalf("node %s has incomplete option: %#v", node.ID, option)
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

	university := findNode(nodes, "university_tier")
	if university == nil {
		t.Fatal("university_tier node not found")
	}
	labels := strings.Join(optionLabels(university.Options), "\n")
	for _, want := range []string{"清北/Top2", "华五", "C9", "985", "211", "双非", "海外 QS100", "海外非 QS100"} {
		if !strings.Contains(labels, want) {
			t.Fatalf("university options = %q, want %q", labels, want)
		}
	}
	if !containsAll(labels, "清北/Top2", "双非") {
		t.Fatalf("university options = %q, want first-degree filter tiers", labels)
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

func optionLabels(options []game.AnswerOption) []string {
	labels := make([]string, 0, len(options))
	for _, option := range options {
		labels = append(labels, option.Label)
	}
	return labels
}
