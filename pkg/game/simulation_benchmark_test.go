package game

import "testing"

func BenchmarkSimulationRun(b *testing.B) {
	nodes := []Node{
		{ID: "gaokao_score", Title: "高考成绩 Benchmark", Unlocks: []string{"province_rank"}, Effects: Effects{BenchScore: 10}},
		{ID: "province_rank", Title: "省排名精度校准", Unlocks: []string{"university_tier"}, Effects: Effects{Anxiety: 8}},
		{ID: "university_tier", Title: "大学档次 Benchmark", Unlocks: []string{"internship_gate"}, Effects: Effects{BenchScore: 12}},
		{ID: "internship_gate", Title: "第一段实习入口", Unlocks: []string{"big_factory_gate"}, Effects: Effects{PeerComparison: 10}},
		{ID: "big_factory_gate", Title: "BigFactory 入口", Unlocks: []string{"ai_replacement_audit"}, Effects: Effects{BenchScore: 18, Anxiety: 16}},
		{ID: "ai_replacement_audit", Title: "AI 替代性审计", Effects: Effects{Absurdity: 20, EscapeIndex: 15}},
	}
	engine := NewEngine(nodes)
	actions := []Action{ActionOptimizeMetric, ActionComparePeers, ActionJobTrack, ActionRefuseMetric}

	for i := 0; i < b.N; i++ {
		_ = RunSimulation(engine, int64(i), actions)
	}
}
