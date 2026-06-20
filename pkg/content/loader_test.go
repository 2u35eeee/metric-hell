package content

import "testing"

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
