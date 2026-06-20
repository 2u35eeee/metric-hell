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

func TestStepSubmissionBucketsGaokaoScoreAndRecordsAudit(t *testing.T) {
	nodes := []Node{
		{
			ID:    InitialNodeID,
			Title: "高考成绩 Benchmark",
			Stage: "高中",
			Input: InputSpec{
				Type:        InputTypeNumber,
				Prompt:      "你的高考分数是多少？",
				Placeholder: "0-750",
			},
			Options: []AnswerOption{
				{
					ID:      "score_700_plus",
					Label:   "700+",
					Min:     floatPtr(700),
					Verdict: "你要是 700 以上还认识这个开发者？系统怀疑样本来源异常。",
					Proof:   "系统把你标记为罕见样本，同时继续要求更多可比较字段。",
					Effects: Effects{BenchScore: 20, Anxiety: 8, PeerComparison: 12},
					Unlocks: []string{"university_tier"},
				},
			},
		},
		{
			ID:    "university_tier",
			Title: "大学层次过滤器",
			Stage: "大学",
		},
	}
	engine := NewEngine(nodes)
	initial := engine.InitialResult(1)

	result, err := engine.StepSubmission(initial.State, Submission{
		NodeID:       InitialNodeID,
		NumericValue: floatPtr(701),
	})
	if err != nil {
		t.Fatalf("StepSubmission returned error: %v", err)
	}
	if result.AuditRecord == nil {
		t.Fatal("AuditRecord = nil, want record")
	}
	if result.AuditRecord.Verdict != "你要是 700 以上还认识这个开发者？系统怀疑样本来源异常。" {
		t.Fatalf("verdict = %q", result.AuditRecord.Verdict)
	}
	if result.AuditRecord.SubmittedLabel != "701" {
		t.Fatalf("submitted label = %q, want 701", result.AuditRecord.SubmittedLabel)
	}
	if len(result.State.AuditTrail) != 1 {
		t.Fatalf("len(audit trail) = %d, want 1", len(result.State.AuditTrail))
	}
	if result.CurrentNode == nil || result.CurrentNode.ID != "university_tier" {
		t.Fatalf("current node = %#v, want university_tier", result.CurrentNode)
	}
}

func TestStepSubmissionSelectsUniversityTierVerdict(t *testing.T) {
	nodes := []Node{
		{
			ID:    InitialNodeID,
			Title: "高考成绩 Benchmark",
			Stage: "高中",
			Input: InputSpec{Type: InputTypeSelect, Prompt: "跳过到大学层次"},
			Options: []AnswerOption{
				{ID: "next", Label: "下一步", Verdict: "继续比较。", Proof: "系统需要第一学历字段。", Unlocks: []string{"university_tier"}},
			},
		},
		{
			ID:    "university_tier",
			Title: "大学层次过滤器",
			Stage: "大学",
			Input: InputSpec{Type: InputTypeSelect, Prompt: "你的第一学历层次被系统归到哪一栏？"},
			Options: []AnswerOption{
				{
					ID:      "top2",
					Label:   "清北/Top2",
					Verdict: "系统记录：Top2 字段点亮。补充提示：还不是分不清鹅腿和鸭腿？",
					Proof:   "第一学历过滤器暂时收起白眼，但马上请求 GPA、实习和厂牌继续证明你。",
					Effects: Effects{BenchScore: 18, PeerComparison: 10},
					Unlocks: []string{"gpa_loop"},
				},
				{
					ID:      "double_non",
					Label:   "双非",
					Verdict: "系统启动第一学历滤镜：不是你不行，是这张表懒得读第二行。",
					Proof:   "系统把复杂经历压成入口字段，并要求你用更多证据补交解释权。",
					Effects: Effects{Anxiety: 12, Selfhood: -4, EscapeIndex: 4},
					Unlocks: []string{"gpa_loop"},
				},
			},
		},
		{ID: "gpa_loop", Title: "GPA 字段维护", Stage: "大学"},
	}
	engine := NewEngine(nodes)
	first, err := engine.StepSubmission(engine.InitialResult(1).State, Submission{NodeID: InitialNodeID, OptionID: "next"})
	if err != nil {
		t.Fatalf("first StepSubmission returned error: %v", err)
	}

	result, err := engine.StepSubmission(first.State, Submission{NodeID: "university_tier", OptionID: "top2"})
	if err != nil {
		t.Fatalf("StepSubmission returned error: %v", err)
	}
	if result.AuditRecord == nil || !strings.Contains(result.AuditRecord.Verdict, "Top2 字段点亮") {
		t.Fatalf("audit record = %#v, want Top2 verdict", result.AuditRecord)
	}
	if result.State.AuditTrail[1].SubmittedLabel != "清北/Top2" {
		t.Fatalf("submitted label = %q, want 清北/Top2", result.State.AuditTrail[1].SubmittedLabel)
	}
}

func TestEngineDoesNotEndAtBigFactoryBeforePositionTier(t *testing.T) {
	nodes := []Node{
		{
			ID:      "big_factory_gate",
			Title:   "BigFactory 入口",
			Stage:   "大厂",
			Input:   InputSpec{Type: InputTypeSelect, Prompt: "是否进入 BigFactory？"},
			Options: []AnswerOption{{ID: "yes", Label: "已进入", Verdict: "BigFactory 字段点亮。", Proof: "系统继续比较岗位。", Effects: Effects{BenchScore: 14, Anxiety: 12}, Unlocks: []string{"position_tier"}}},
		},
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

	result, err := engine.StepSubmission(state, Submission{NodeID: "big_factory_gate", OptionID: "yes"})
	if err != nil {
		t.Fatalf("StepSubmission returned error: %v", err)
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
		{
			ID:      "gaokao_score",
			Title:   "高考成绩 Benchmark",
			Stage:   "高中",
			Input:   InputSpec{Type: InputTypeSelect, Prompt: "提交字段"},
			Options: []AnswerOption{{ID: "next", Label: "继续", Verdict: "继续比较。", Proof: "系统继续生成字段。", Effects: Effects{BenchScore: 20}, Unlocks: []string{"province_rank"}}},
		},
		{
			ID:      "province_rank",
			Title:   "省排名精度校准",
			Stage:   "高中",
			Input:   InputSpec{Type: InputTypeSelect, Prompt: "提交字段"},
			Options: []AnswerOption{{ID: "next", Label: "继续", Verdict: "继续比较。", Proof: "系统继续生成字段。", Effects: Effects{Anxiety: 15}, Unlocks: []string{"life_not_workflow"}}},
		},
		{
			ID:      "life_not_workflow",
			Title:   "人生不是 Workflow",
			Stage:   "逃逸",
			Input:   InputSpec{Type: InputTypeSelect, Prompt: "结束"},
			Options: []AnswerOption{{ID: "end", Label: "暂停", Verdict: "系统停止。", Proof: "你离开了表格。", Effects: Effects{Selfhood: 30, EscapeIndex: 35}}},
		},
	})

	result := RunSimulation(engine, 42, nil)

	if result.Ending == nil {
		t.Fatal("Ending = nil, want an ending")
	}
}

func floatPtr(value float64) *float64 {
	return &value
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
