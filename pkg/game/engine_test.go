package game

import (
	"strings"
	"testing"
)

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

func TestNodeBranchesProvideContextualActionsAndUnlockChosenPath(t *testing.T) {
	nodes := []Node{
		{
			ID:        "gaokao_score",
			Title:     "高考成绩 Benchmark",
			Stage:     "高中",
			Questions: []string{"分数之后，系统还需要什么字段？"},
			Branches: []BranchOption{
				{
					ID:          Action("submit_rank_slice"),
					Label:       "把分数切成更细排名",
					Scene:       "家族群正在等待一个可以截图传播的数字。",
					Description: "选择后系统会进入省排名精度校准。",
					Effects:     Effects{BenchScore: 4, Anxiety: 5, PeerComparison: 8},
					Unlocks:     []string{"province_rank"},
					ResultText:  "你提交了更细排名。系统满意地增加了比较分辨率。",
				},
				{
					ID:          Action("protect_curiosity"),
					Label:       "先保留一点好奇心",
					Scene:       "你没有继续把自己拆成小数点后的排名。",
					Description: "选择后系统会提前暴露逃逸路径。",
					Effects:     Effects{Selfhood: 8, Curiosity: 7, EscapeIndex: 10},
					Unlocks:     []string{"life_not_workflow"},
					ResultText:  "你保留了一个不便排序的部分。系统开始发热。",
				},
			},
			Effects: Effects{Absurdity: 2},
		},
		{ID: "province_rank", Title: "省排名精度校准", Stage: "高中"},
		{ID: "life_not_workflow", Title: "人生不是 Workflow", Stage: "逃逸"},
	}
	engine := NewEngine(nodes)
	initial := engine.InitialResult(7)

	if len(initial.Actions) != 2 {
		t.Fatalf("len(actions) = %d, want 2", len(initial.Actions))
	}
	if initial.Actions[0].Scene == "" {
		t.Fatalf("first action scene is empty: %#v", initial.Actions[0])
	}

	result, err := engine.Step(initial.State, Action("protect_curiosity"))
	if err != nil {
		t.Fatalf("Step returned error: %v", err)
	}
	if Contains(result.State.UnlockedNodes, "province_rank") {
		t.Fatalf("province_rank unlocked after choosing escape branch: %#v", result.State.UnlockedNodes)
	}
	if !Contains(result.State.UnlockedNodes, "life_not_workflow") {
		t.Fatalf("life_not_workflow not unlocked: %#v", result.State.UnlockedNodes)
	}
	if result.CurrentNode == nil || result.CurrentNode.ID != "life_not_workflow" {
		t.Fatalf("current node = %#v, want life_not_workflow", result.CurrentNode)
	}
	if result.State.Curiosity <= initial.State.Curiosity {
		t.Fatalf("curiosity = %d, want above %d", result.State.Curiosity, initial.State.Curiosity)
	}
	latestLog := strings.Join(result.State.EventLog, "\n")
	if !strings.Contains(latestLog, "人生不是 Workflow") {
		t.Fatalf("event log = %q, want unlocked node title", latestLog)
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
