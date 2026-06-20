package game

import "testing"

func TestSpawnNewBenchmarksDoesNotDuplicateUnlockedNodes(t *testing.T) {
	state := State{UnlockedNodes: []string{"province_rank"}}
	node := Node{
		ID:      "gaokao_score",
		Title:   "高考成绩 Benchmark",
		Unlocks: []string{"province_rank", "university_tier"},
	}

	got := SpawnNewBenchmarks(state, node)

	if countOccurrences(got.UnlockedNodes, "province_rank") != 1 {
		t.Fatalf("province_rank duplicated in %#v", got.UnlockedNodes)
	}
	if !Contains(got.UnlockedNodes, "university_tier") {
		t.Fatalf("university_tier not unlocked: %#v", got.UnlockedNodes)
	}
}

func TestEngineStepCompletesNodeAndReturnsNextNode(t *testing.T) {
	nodes := []Node{
		{
			ID:        "gaokao_score",
			Title:     "高考成绩 Benchmark",
			Stage:     "高中",
			Questions: []string{"分数之后，系统还需要什么字段？"},
			Unlocks:   []string{"province_rank"},
			Effects:   Effects{BenchScore: 5, Anxiety: 3},
		},
		{ID: "province_rank", Title: "省排名精度校准", Stage: "高中"},
	}
	engine := NewEngine(nodes)
	state := NewInitialState(1)

	result, err := engine.Step(state, ActionOptimizeMetric)
	if err != nil {
		t.Fatalf("Step returned error: %v", err)
	}
	if !Contains(result.State.CompletedNodes, "gaokao_score") {
		t.Fatalf("completed nodes = %#v, want gaokao_score", result.State.CompletedNodes)
	}
	if result.CurrentNode == nil || result.CurrentNode.ID != "province_rank" {
		t.Fatalf("current node = %#v, want province_rank", result.CurrentNode)
	}
}

func TestEngineDoesNotEndAtBigFactoryBeforePositionTier(t *testing.T) {
	nodes := []Node{
		{ID: "big_factory_gate", Title: "BigFactory 入口", Stage: "大厂", Unlocks: []string{"position_tier"}, Effects: Effects{BenchScore: 14, Anxiety: 12}},
		{ID: "position_tier", Title: "大厂岗位档次 Benchmark", Stage: "大厂"},
	}
	engine := NewEngine(nodes)
	state := State{
		Stage:         "大厂",
		BenchScore:    80,
		Anxiety:       80,
		Selfhood:      50,
		Energy:        50,
		UnlockedNodes: []string{"big_factory_gate"},
	}

	result, err := engine.Step(state, ActionOptimizeMetric)
	if err != nil {
		t.Fatalf("Step returned error: %v", err)
	}
	if result.Ended {
		t.Fatalf("Ended = true, want false so position_tier can be played")
	}
	if result.CurrentNode == nil || result.CurrentNode.ID != "position_tier" {
		t.Fatalf("current node = %#v, want position_tier", result.CurrentNode)
	}
}

func TestBenchmarkSimulationRunCompletesWithEnding(t *testing.T) {
	engine := NewEngine([]Node{
		{ID: "gaokao_score", Title: "高考成绩 Benchmark", Stage: "高中", Unlocks: []string{"province_rank"}, Effects: Effects{BenchScore: 20}},
		{ID: "province_rank", Title: "省排名精度校准", Stage: "高中", Unlocks: []string{"life_not_workflow"}, Effects: Effects{Anxiety: 15}},
		{ID: "life_not_workflow", Title: "人生不是 Workflow", Stage: "逃逸", Effects: Effects{Selfhood: 30, EscapeIndex: 35}},
	})

	result := RunSimulation(engine, 42, []Action{ActionOptimizeMetric, ActionRefuseMetric, ActionRefuseMetric})

	if result.Ending == nil {
		t.Fatal("Ending = nil, want an ending")
	}
}

func countOccurrences(values []string, needle string) int {
	count := 0
	for _, value := range values {
		if value == needle {
			count++
		}
	}
	return count
}
