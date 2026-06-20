package game

import "testing"

func BenchmarkSimulationRun(b *testing.B) {
	nodes := []Node{
		benchNode("gaokao_score", "高考成绩 Benchmark", "province_rank", Effects{BenchScore: 10}),
		benchNode("province_rank", "省排名精度校准", "university_tier", Effects{Anxiety: 8}),
		benchNode("university_tier", "大学档次 Benchmark", "internship_gate", Effects{BenchScore: 12}),
		benchNode("internship_gate", "第一段实习入口", "big_factory_gate", Effects{PeerComparison: 10}),
		benchNode("big_factory_gate", "BigFactory 入口", "ai_replacement_audit", Effects{BenchScore: 18, Anxiety: 16}),
		benchNode("ai_replacement_audit", "AI 替代性审计", "", Effects{Absurdity: 20, EscapeIndex: 15}),
	}
	engine := NewEngine(nodes)

	for i := 0; i < b.N; i++ {
		_ = RunSimulation(engine, int64(i), nil)
	}
}

func benchNode(id string, title string, unlock string, effects Effects) Node {
	unlocks := []string{}
	if unlock != "" {
		unlocks = []string{unlock}
	}
	return Node{
		ID:    id,
		Title: title,
		Input: InputSpec{Type: InputTypeSelect, Prompt: "提交字段"},
		Options: []AnswerOption{
			{ID: "next", Label: "继续", Verdict: "继续比较。", Proof: "系统继续生成字段。", Effects: effects, Unlocks: unlocks},
		},
	}
}
