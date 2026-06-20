package game

import "testing"

func TestEvaluateEndingReturnsEscapeEnding(t *testing.T) {
	state := State{EscapeIndex: 80, Selfhood: 70}

	ending := EvaluateEnding(state)

	if ending == nil {
		t.Fatal("ending = nil, want escape ending")
	}
	if ending.ID != "escape_success" {
		t.Fatalf("ending ID = %q, want escape_success", ending.ID)
	}
}

func TestEvaluateEndingReturnsBigFactoryEnding(t *testing.T) {
	state := State{
		BenchScore:     75,
		Anxiety:        78,
		CompletedNodes: []string{"big_factory_gate", "position_tier"},
	}

	ending := EvaluateEnding(state)

	if ending == nil {
		t.Fatal("ending = nil, want big_factory_metric")
	}
	if ending.ID != "big_factory_metric" {
		t.Fatalf("ending ID = %q, want big_factory_metric", ending.ID)
	}
}
