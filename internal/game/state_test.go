package game

import "testing"

func TestClampStateKeepsCoreValuesWithinRange(t *testing.T) {
	s := State{
		BenchScore:     120,
		Anxiety:        -5,
		Selfhood:       150,
		Energy:         -20,
		Curiosity:      101,
		ParentPressure: -1,
		PeerComparison: 200,
		EscapeIndex:    -30,
		Absurdity:      999,
	}

	got := ClampState(s)

	for name, value := range map[string]int{
		"BenchScore":     got.BenchScore,
		"Anxiety":        got.Anxiety,
		"Selfhood":       got.Selfhood,
		"Energy":         got.Energy,
		"Curiosity":      got.Curiosity,
		"ParentPressure": got.ParentPressure,
		"PeerComparison": got.PeerComparison,
		"EscapeIndex":    got.EscapeIndex,
		"Absurdity":      got.Absurdity,
	} {
		if value < 0 || value > 100 {
			t.Fatalf("%s = %d, want 0..100", name, value)
		}
	}
}
