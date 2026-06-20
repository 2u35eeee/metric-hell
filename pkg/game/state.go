package game

import "fmt"

const InitialNodeID = "gaokao_score"

type State struct {
	Age              int      `json:"age"`
	Stage            string   `json:"stage"`
	BenchScore       int      `json:"bench_score"`
	Anxiety          int      `json:"anxiety"`
	Selfhood         int      `json:"selfhood"`
	Energy           int      `json:"energy"`
	Curiosity        int      `json:"curiosity"`
	ParentPressure   int      `json:"parent_pressure"`
	PeerComparison   int      `json:"peer_comparison"`
	EscapeIndex      int      `json:"escape_index"`
	Absurdity        int      `json:"absurdity"`
	CompletedNodes   []string `json:"completed_nodes"`
	UnlockedNodes    []string `json:"unlocked_nodes"`
	EventLog         []string `json:"event_log"`
	Turn             int      `json:"turn"`
	VirtualStudentID string   `json:"virtual_student_id"`
}

func NewInitialState(seed int64) State {
	suffix := int(seed%1000+1000) % 1000
	return State{
		Age:            18,
		Stage:          "高中",
		BenchScore:     18,
		Anxiety:        32,
		Selfhood:       72,
		Energy:         76,
		Curiosity:      68,
		ParentPressure: 45,
		PeerComparison: 30,
		EscapeIndex:    12,
		Absurdity:      10,
		UnlockedNodes:  []string{InitialNodeID},
		CompletedNodes: []string{},
		EventLog: []string{
			fmt.Sprintf("虚构学生 #%03d 已生成。系统声明：该角色不对应任何真实个人。", suffix),
			"WorkflowBench 已启动：请提交一个可以被排序的人生片段。",
		},
		Turn:             0,
		VirtualStudentID: fmt.Sprintf("虚构学生 #%03d", suffix),
	}
}

func ClampState(s State) State {
	s.BenchScore = clamp(s.BenchScore)
	s.Anxiety = clamp(s.Anxiety)
	s.Selfhood = clamp(s.Selfhood)
	s.Energy = clamp(s.Energy)
	s.Curiosity = clamp(s.Curiosity)
	s.ParentPressure = clamp(s.ParentPressure)
	s.PeerComparison = clamp(s.PeerComparison)
	s.EscapeIndex = clamp(s.EscapeIndex)
	s.Absurdity = clamp(s.Absurdity)
	return s
}

func clamp(value int) int {
	if value < 0 {
		return 0
	}
	if value > 100 {
		return 100
	}
	return value
}

func Contains(values []string, needle string) bool {
	for _, value := range values {
		if value == needle {
			return true
		}
	}
	return false
}

func appendUnique(values []string, value string) []string {
	if Contains(values, value) {
		return values
	}
	return append(values, value)
}
