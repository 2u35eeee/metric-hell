package game

type Action string

const (
	ActionOptimizeMetric Action = "optimize_metric"
	ActionComparePeers   Action = "compare_peers"
	ActionRest           Action = "rest"
	ActionSwitchTrack    Action = "switch_track"
	ActionJobTrack       Action = "job_track"
	ActionStableTrack    Action = "stable_track"
	ActionRefuseMetric   Action = "refuse_metric"
)

type ActionOption struct {
	ID          Action   `json:"id"`
	Label       string   `json:"label"`
	Scene       string   `json:"scene,omitempty"`
	Description string   `json:"description"`
	Effects     Effects  `json:"effects,omitempty"`
	Unlocks     []string `json:"unlocks,omitempty"`
}

var defaultActions = []Action{
	ActionOptimizeMetric,
	ActionComparePeers,
	ActionRest,
	ActionSwitchTrack,
	ActionJobTrack,
	ActionStableTrack,
	ActionRefuseMetric,
}

func AvailableActions(_ State, node *Node) []ActionOption {
	if node != nil && len(node.Branches) > 0 {
		options := make([]ActionOption, 0, len(node.Branches))
		for _, branch := range node.Branches {
			options = append(options, ActionOption{
				ID:          branch.ID,
				Label:       branch.Label,
				Scene:       branch.Scene,
				Description: branch.Description,
				Effects:     branch.Effects,
				Unlocks:     branch.Unlocks,
			})
		}
		return options
	}

	options := make([]ActionOption, 0, len(defaultActions))
	for _, action := range defaultActions {
		options = append(options, ActionOption{
			ID:          action,
			Label:       ActionLabel(action),
			Description: ActionDescription(action),
		})
	}
	return options
}

func FindBranch(node *Node, action Action) (BranchOption, bool) {
	if node == nil {
		return BranchOption{}, false
	}
	for _, branch := range node.Branches {
		if branch.ID == action {
			return branch, true
		}
	}
	return BranchOption{}, false
}

func ApplyAction(s State, action Action) State {
	switch action {
	case ActionOptimizeMetric:
		s.BenchScore += 10
		s.Anxiety += 8
		s.Energy -= 10
		s.Selfhood -= 4
		s.EventLog = append(s.EventLog, "你选择继续优化指标。系统记录：排序字段更整齐了，自我叙事略微变薄。")
	case ActionComparePeers:
		s.PeerComparison += 15
		s.Anxiety += 12
		s.Energy -= 5
		s.Absurdity += 5
		s.EventLog = append(s.EventLog, "你打开同辈比较面板。系统提示：样本量不足，但焦虑量充足。")
	case ActionRest:
		s.Energy += 12
		s.Anxiety -= 8
		s.BenchScore -= 3
		s.ParentPressure += 3
		s.EventLog = append(s.EventLog, "你选择休息。系统暂时无法理解该行为的 KPI 价值。")
	case ActionSwitchTrack:
		s.Curiosity += 12
		s.EscapeIndex += 8
		s.Anxiety += 5
		s.BenchScore -= 5
		s.EventLog = append(s.EventLog, "你尝试换赛道。系统正在把非标准路径转换为新的表格字段。")
	case ActionJobTrack:
		s.BenchScore += 8
		s.Anxiety += 8
		s.PeerComparison += 6
		s.Energy -= 6
		s.EventLog = append(s.EventLog, "你进入就业路线。系统开始加载厂牌、岗位、转正率和裁员风险。")
	case ActionStableTrack:
		s.ParentPressure -= 8
		s.BenchScore += 5
		s.Anxiety += 10
		s.Curiosity -= 5
		s.EventLog = append(s.EventLog, "你选择稳定路线。外部风险下降，新的稳定性比较维度已生成。")
	case ActionRefuseMetric:
		s.Selfhood += 15
		s.EscapeIndex += 12
		s.ParentPressure += 8
		s.BenchScore -= 8
		s.Absurdity += 10
		s.EventLog = append(s.EventLog, "你拒绝被单一指标解释。系统报错：该输入无法稳定归一化。")
	default:
		s.EventLog = append(s.EventLog, "系统收到未知动作，已将其归档为：非标准行为。")
	}
	return ClampState(s)
}

func ActionLabel(action Action) string {
	switch action {
	case ActionOptimizeMetric:
		return "继续优化指标"
	case ActionComparePeers:
		return "打开同辈比较面板"
	case ActionRest:
		return "休息一下"
	case ActionSwitchTrack:
		return "换一条赛道"
	case ActionJobTrack:
		return "进入就业路线"
	case ActionStableTrack:
		return "进入稳定路线"
	case ActionRefuseMetric:
		return "拒绝被指标解释"
	default:
		return "未知动作"
	}
}

func ActionDescription(action Action) string {
	switch action {
	case ActionOptimizeMetric:
		return "BenchScore 上升，但焦虑和疲惫也会增加。"
	case ActionComparePeers:
		return "比较维度变多，系统显得更开心，人不一定。"
	case ActionRest:
		return "恢复精力，不等于失败；只是系统会短暂困惑。"
	case ActionSwitchTrack:
		return "提升好奇心和逃逸指数，但会暂时降低可排序性。"
	case ActionJobTrack:
		return "进入厂牌、岗位和转正率构成的新表格。"
	case ActionStableTrack:
		return "降低部分外部压力，同时生成稳定性比较字段。"
	case ActionRefuseMetric:
		return "提升自我感和逃逸指数，系统会尝试继续解释你。"
	default:
		return "系统尚未理解该动作。"
	}
}
