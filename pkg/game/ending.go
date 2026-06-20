package game

type Ending struct {
	ID               string `json:"id"`
	Title            string `json:"title"`
	Type             string `json:"type"`
	SystemEvaluation string `json:"system_evaluation"`
	HiddenEvaluation string `json:"hidden_evaluation"`
}

func EvaluateEnding(s State) *Ending {
	switch {
	case s.EscapeIndex >= 75 && s.Selfhood >= 65:
		return &Ending{
			ID:               "escape_success",
			Title:            "指标逃逸型",
			Type:             "人生不是 Workflow",
			SystemEvaluation: "该玩家路径异常，无法归类为单一字段。",
			HiddenEvaluation: "你没有完成所有 benchmark，但你重新获得了用自己的语言描述自己的能力。",
		}
	case s.BenchScore >= 85 && s.Selfhood <= 30 && s.Anxiety >= 75:
		return &Ending{
			ID:               "perfect_hollow",
			Title:            "高 Bench 低自我型",
			Type:             "满分空心人",
			SystemEvaluation: "你非常适合被表格展示。",
			HiddenEvaluation: "你几乎完成了所有可见指标，但系统仍无法回答：你喜欢什么？",
		}
	case Contains(s.CompletedNodes, "big_factory_gate") && s.BenchScore >= 70 && s.Anxiety >= 70:
		return &Ending{
			ID:               "big_factory_metric",
			Title:            "BigFactory 指标人",
			Type:             "大厂指标人",
			SystemEvaluation: "你已进入 BigFactory。系统解锁岗位层级、转正率和风险巡检。",
			HiddenEvaluation: "benchmark 并没有结束，只是换了一套更贵的 UI。",
		}
	case Contains(s.CompletedNodes, "layoff_risk_check") && s.Anxiety >= 65:
		return &Ending{
			ID:               "layoff_dashboard",
			Title:            "风险面板常驻型",
			Type:             "裁员风险巡检",
			SystemEvaluation: "岗位状态：已获得。风险状态：持续刷新。",
			HiddenEvaluation: "系统无法承诺安全感，只能继续生产概率。",
		}
	case Contains(s.CompletedNodes, "ai_replacement_audit") && s.Absurdity >= 65:
		return &Ending{
			ID:               "ai_replacement_loop",
			Title:            "替代性审计型",
			Type:             "AI 替代性审计",
			SystemEvaluation: "你刚刚适配岗位，岗位定义已经开始移动。",
			HiddenEvaluation: "系统把变化解释成风险，却忘了变化也可能是重新选择的机会。",
		}
	case s.Turn >= 20 || Contains(s.CompletedNodes, "life_not_workflow"):
		return &Ending{
			ID:               "workflow_overflow",
			Title:            "指标继续生成中",
			Type:             "人生完成度无法达到 100%",
			SystemEvaluation: "系统发现新的评价维度，通关条件已自动延期。",
			HiddenEvaluation: "你没有失败。只是这个版本的评价器不支持完整的人。",
		}
	default:
		return nil
	}
}
