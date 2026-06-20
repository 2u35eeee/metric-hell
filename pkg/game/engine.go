package game

import (
	"errors"
	"fmt"
	"strconv"
)

type Engine struct {
	nodes []Node
	byID  map[string]Node
}

type Result struct {
	State       State        `json:"state"`
	CurrentNode *Node        `json:"current_node"`
	AuditRecord *AuditRecord `json:"audit_record,omitempty"`
	Ended       bool         `json:"ended"`
	Ending      *Ending      `json:"ending"`
}

type Submission struct {
	NodeID       string   `json:"node_id"`
	NumericValue *float64 `json:"numeric_value,omitempty"`
	OptionID     string   `json:"option_id,omitempty"`
}

func NewEngine(nodes []Node) *Engine {
	byID := make(map[string]Node, len(nodes))
	for _, node := range nodes {
		byID[node.ID] = node
	}
	return &Engine{nodes: nodes, byID: byID}
}

func (e *Engine) InitialResult(seed int64) Result {
	state := NewInitialState(seed)
	node := e.NextNode(state)
	return Result{
		State:       state,
		CurrentNode: node,
		Ended:       false,
		Ending:      nil,
	}
}

func (e *Engine) StepSubmission(state State, submission Submission) (Result, error) {
	if len(e.nodes) == 0 {
		return Result{}, errors.New("no benchmark nodes loaded")
	}

	current := e.NextNode(state)
	if current == nil {
		ending := EvaluateEnding(state)
		if ending == nil {
			ending = &Ending{
				ID:               "no_more_nodes",
				Title:            "指标耗尽异常",
				Type:             "系统边界条件",
				SystemEvaluation: "系统暂时没有更多字段可以生成。",
				HiddenEvaluation: "这可能是今天最接近自由的一刻。",
			}
		}
		return Result{State: state, Ended: true, Ending: ending}, nil
	}
	if submission.NodeID != "" && submission.NodeID != current.ID {
		return Result{}, fmt.Errorf("submission node %q does not match current node %q", submission.NodeID, current.ID)
	}

	option, submittedLabel, err := matchSubmission(*current, submission)
	if err != nil {
		return Result{}, err
	}

	unlocks := option.Unlocks
	if unlocks == nil {
		unlocks = current.Unlocks
	}

	state = ApplyEffects(state, option.Effects)
	state = ApplyEffects(state, current.Effects)
	state.Absurdity += current.Absurdity
	state.Stage = current.Stage
	state.Turn++
	state.CompletedNodes = appendUnique(state.CompletedNodes, current.ID)

	record := AuditRecord{
		Turn:           state.Turn,
		NodeID:         current.ID,
		NodeTitle:      current.Title,
		Stage:          current.Stage,
		Prompt:         current.Input.Prompt,
		SubmittedLabel: submittedLabel,
		OptionID:       option.ID,
		Verdict:        option.Verdict,
		Proof:          option.Proof,
		Effects:        option.Effects,
		Unlocks:        append([]string(nil), unlocks...),
	}
	state.AuditTrail = append(state.AuditTrail, record)
	state.EventLog = append(state.EventLog, option.Verdict)
	state = e.spawnNewBenchmarks(state, *current, unlocks)
	state = ClampState(state)

	next := e.NextNode(state)
	ending := EvaluateEnding(state)
	if ending != nil && (next == nil || isImmediateEnding(ending.ID)) {
		return Result{
			State:       state,
			CurrentNode: nil,
			AuditRecord: &record,
			Ended:       true,
			Ending:      ending,
		}, nil
	}

	return Result{
		State:       state,
		CurrentNode: next,
		AuditRecord: &record,
		Ended:       false,
		Ending:      nil,
	}, nil
}

func matchSubmission(node Node, submission Submission) (AnswerOption, string, error) {
	if len(node.Options) == 0 {
		return AnswerOption{}, "", fmt.Errorf("node %q has no answer options", node.ID)
	}

	switch node.Input.Type {
	case InputTypeNumber:
		if submission.NumericValue == nil {
			return AnswerOption{}, "", fmt.Errorf("node %q requires numeric_value", node.ID)
		}
		value := *submission.NumericValue
		for _, option := range node.Options {
			if option.matchesNumber(value) {
				return option, strconv.FormatFloat(value, 'f', -1, 64), nil
			}
		}
		return AnswerOption{}, "", fmt.Errorf("numeric value %s did not match node %q", strconv.FormatFloat(value, 'f', -1, 64), node.ID)
	case InputTypeSelect:
		if submission.OptionID == "" {
			return AnswerOption{}, "", fmt.Errorf("node %q requires option_id", node.ID)
		}
		for _, option := range node.Options {
			if option.ID == submission.OptionID {
				return option, option.Label, nil
			}
		}
		return AnswerOption{}, "", fmt.Errorf("option %q did not match node %q", submission.OptionID, node.ID)
	default:
		return AnswerOption{}, "", fmt.Errorf("node %q has unsupported input type %q", node.ID, node.Input.Type)
	}
}

func (option AnswerOption) matchesNumber(value float64) bool {
	if option.Min != nil && value < *option.Min {
		return false
	}
	if option.Max != nil && value > *option.Max {
		return false
	}
	return true
}

func (e *Engine) NextNode(state State) *Node {
	for _, id := range state.UnlockedNodes {
		if Contains(state.CompletedNodes, id) {
			continue
		}
		node, ok := e.byID[id]
		if !ok {
			continue
		}
		return &node
	}
	return nil
}

func SpawnNewBenchmarks(s State, node Node) State {
	return SpawnNewBenchmarksWithUnlocks(s, node, node.Unlocks)
}

func SpawnNewBenchmarksWithUnlocks(s State, node Node, unlocks []string) State {
	for _, next := range unlocks {
		if Contains(s.UnlockedNodes, next) {
			continue
		}
		s.UnlockedNodes = append(s.UnlockedNodes, next)
		s.EventLog = append(s.EventLog, fmt.Sprintf("系统检测到你完成了「%s」，已生成更细评价指标：%s。", node.Title, next))
	}
	return s
}

func (e *Engine) spawnNewBenchmarks(s State, node Node, unlocks []string) State {
	for _, next := range unlocks {
		if Contains(s.UnlockedNodes, next) {
			continue
		}
		s.UnlockedNodes = append(s.UnlockedNodes, next)
		label := next
		if nextNode, ok := e.byID[next]; ok {
			label = nextNode.Title
		}
		s.EventLog = append(s.EventLog, fmt.Sprintf("系统检测到你完成了「%s」，已生成更细评价指标：%s。", node.Title, label))
	}
	return s
}

func RunSimulation(engine *Engine, seed int64, submissions []Submission) Result {
	result := engine.InitialResult(seed)
	for i := 0; i < 64 && !result.Ended; i++ {
		if result.CurrentNode == nil {
			break
		}
		submission := defaultSubmission(*result.CurrentNode)
		if len(submissions) > 0 {
			submission = submissions[i%len(submissions)]
			if submission.NodeID == "" {
				submission.NodeID = result.CurrentNode.ID
			}
		}
		next, err := engine.StepSubmission(result.State, submission)
		if err != nil {
			return Result{
				State: result.State,
				Ended: true,
				Ending: &Ending{
					ID:               "simulation_error",
					Title:            "模拟异常",
					Type:             "系统边界条件",
					SystemEvaluation: err.Error(),
					HiddenEvaluation: "评价器自己也不是很稳定。",
				},
			}
		}
		result = next
	}
	return result
}

func defaultSubmission(node Node) Submission {
	submission := Submission{NodeID: node.ID}
	if len(node.Options) == 0 {
		return submission
	}
	first := node.Options[0]
	if node.Input.Type == InputTypeNumber {
		switch {
		case first.Min != nil:
			submission.NumericValue = first.Min
		case first.Max != nil:
			submission.NumericValue = first.Max
		default:
			value := 0.0
			submission.NumericValue = &value
		}
		return submission
	}
	submission.OptionID = first.ID
	return submission
}

func isImmediateEnding(id string) bool {
	return id == "escape_success" || id == "perfect_hollow"
}
