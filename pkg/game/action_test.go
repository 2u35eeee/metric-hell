package game

import "testing"

func TestApplyActionOptimizeMetricTradesSelfhoodForBenchScore(t *testing.T) {
	start := State{
		BenchScore: 40,
		Anxiety:    30,
		Selfhood:   70,
		Energy:     80,
	}

	got := ApplyAction(start, ActionOptimizeMetric)

	if got.BenchScore != 50 {
		t.Fatalf("BenchScore = %d, want 50", got.BenchScore)
	}
	if got.Anxiety != 38 {
		t.Fatalf("Anxiety = %d, want 38", got.Anxiety)
	}
	if got.Selfhood != 66 {
		t.Fatalf("Selfhood = %d, want 66", got.Selfhood)
	}
	if got.Energy != 70 {
		t.Fatalf("Energy = %d, want 70", got.Energy)
	}
}

func TestActionLabelReturnsReadableText(t *testing.T) {
	if got := ActionLabel(ActionRefuseMetric); got == "" || got == string(ActionRefuseMetric) {
		t.Fatalf("ActionLabel(ActionRefuseMetric) = %q, want human-readable label", got)
	}
}
