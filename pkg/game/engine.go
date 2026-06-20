package game

import (
	"errors"
	"fmt"
)

type Engine struct {
	nodes []Node
	byID  map[string]Node
}

type Result struct {
	State       State          `json:"state"`
	CurrentNode *Node          `json:"current_node"`
	Actions     []ActionOption `json:"actions"`
	Ended       bool           `json:"ended"`
	Ending      *Ending        `json:"ending"`
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
		Actions:     AvailableActions(state, node),
		Ended:       false,
		Ending:      nil,
	}
}

func (e *Engine) Step(state State, action Action) (Result, error) {
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

	unlocks := current.Unlocks
	if branch, ok := FindBranch(current, action); ok {
		state = ApplyEffects(state, branch.Effects)
		if branch.ResultText != "" {
			state.EventLog = append(state.EventLog, branch.ResultText)
		} else {
			state.EventLog = append(state.EventLog, fmt.Sprintf("你选择了「%s」。系统正在把该选择转换为可比较字段。", branch.Label))
		}
		if branch.Unlocks != nil {
			unlocks = branch.Unlocks
		}
	} else {
		state = ApplyAction(state, action)
	}
	state = ApplyEffects(state, current.Effects)
	state.Absurdity += current.Absurdity
	state.Stage = current.Stage
	state.Turn++
	state.CompletedNodes = appendUnique(state.CompletedNodes, current.ID)

	if current.TextOnPass != "" {
		state.EventLog = append(state.EventLog, current.TextOnPass)
	} else {
		state.EventLog = append(state.EventLog, fmt.Sprintf("系统检测到你完成了「%s」。正在生成更精确的问题。", current.Title))
	}
	state = e.spawnNewBenchmarks(state, *current, unlocks)
	state = ClampState(state)

	next := e.NextNode(state)
	ending := EvaluateEnding(state)
	if ending != nil && (next == nil || isImmediateEnding(ending.ID)) {
		return Result{
			State:       state,
			CurrentNode: nil,
			Actions:     nil,
			Ended:       true,
			Ending:      ending,
		}, nil
	}

	return Result{
		State:       state,
		CurrentNode: next,
		Actions:     AvailableActions(state, next),
		Ended:       false,
		Ending:      nil,
	}, nil
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

func RunSimulation(engine *Engine, seed int64, actions []Action) Result {
	result := engine.InitialResult(seed)
	if len(actions) == 0 {
		actions = defaultActions
	}
	for i := 0; i < 64 && !result.Ended; i++ {
		action := actions[i%len(actions)]
		next, err := engine.Step(result.State, action)
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

func isImmediateEnding(id string) bool {
	return id == "escape_success" || id == "perfect_hollow"
}
