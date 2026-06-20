package game

type Node struct {
	ID          string   `json:"id"`
	Title       string   `json:"title"`
	Stage       string   `json:"stage"`
	Questions   []string `json:"questions"`
	Unlocks     []string `json:"unlocks"`
	TextOnEnter string   `json:"text_on_enter"`
	TextOnPass  string   `json:"text_on_pass"`
	TextOnFail  string   `json:"text_on_fail"`
	Effects     Effects  `json:"effects"`
	Absurdity   int      `json:"absurdity"`
}

type Effects struct {
	BenchScore     int `json:"bench_score"`
	Anxiety        int `json:"anxiety"`
	Selfhood       int `json:"selfhood"`
	Energy         int `json:"energy"`
	Curiosity      int `json:"curiosity"`
	ParentPressure int `json:"parent_pressure"`
	PeerComparison int `json:"peer_comparison"`
	EscapeIndex    int `json:"escape_index"`
	Absurdity      int `json:"absurdity"`
}

func ApplyEffects(s State, effects Effects) State {
	s.BenchScore += effects.BenchScore
	s.Anxiety += effects.Anxiety
	s.Selfhood += effects.Selfhood
	s.Energy += effects.Energy
	s.Curiosity += effects.Curiosity
	s.ParentPressure += effects.ParentPressure
	s.PeerComparison += effects.PeerComparison
	s.EscapeIndex += effects.EscapeIndex
	s.Absurdity += effects.Absurdity
	return ClampState(s)
}
