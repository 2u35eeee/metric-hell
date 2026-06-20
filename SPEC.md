# WorkflowBench / Metric Hell MVP Spec

## 目标

实现一个可上线试玩的 Go 前后端小游戏。玩家打开 URL 后，通过下拉框选择每回合行动，体验 benchmark 不断增殖的荒诞感。

## 主题

人生不是 workflow，但评价系统总想把它 workflow 化。

游戏从高考成绩和省排名开始，进入大学档次、专业转译、出国档次、GPA、竞赛权重、简历关键词、实习厂牌、实习岗位、大厂岗位、裁员风险和 AI 替代性审计。

## 非目标

- 不预测真实人生。
- 不提供升学、就业或职业建议。
- 不评价真实学校、地区、专业、职业、公司或学生。
- 不做登录、排行榜、数据库、结果存储或分析埋点。

## 架构

- `cmd/server`：本地 Go HTTP 服务。
- `api/index.go`：Vercel Go Runtime 入口。
- `pkg/game`：核心规则引擎。
- `pkg/content`：JSON 节点加载。
- `pkg/api`：HTTP handler。
- `data/nodes.json`：benchmark 内容配置。
- `web`：静态前端。

## 玩法结构

每个 benchmark 节点不是单一路线检查点，而是一个小型人生分叉：

- `scenario` 描述玩家当前处在什么具体场景里。
- `measurement` 解释该指标到底在讽刺什么、范围/口径是什么。
- `questions` 是系统追问，用来制造“评价继续增殖”的荒诞感。
- `branches` 是本回合可选分支，每个分支必须说明：
  - `label`：玩家看到的选择名称。
  - `scene`：这个选择在当前场景里意味着什么。
  - `description`：选择后会发生什么。
  - `effects`：对 0-100 指标面板的影响。
  - `unlocks`：该分支进入的后续节点。
  - `result_text`：选择后的事件日志。

前端必须把分支展示成可读的选择卡，而不是只给一个抽象下拉框。
指标面板必须展示中文指标名、当前值 `/100` 和讽刺解释，不应只泄露内部字段名。

## API

### POST `/api/new`

创建一局虚构学生状态。

### POST `/api/action`

请求：

```json
{
  "state": {},
  "action": "optimize_metric"
}
```

响应：

```json
{
  "state": {},
  "current_node": {
    "scenario": "当前人生场景",
    "measurement": "该指标的口径/范围说明"
  },
  "actions": [
    {
      "id": "gpa_decimal_grind",
      "label": "把 GPA 卷到小数点后一位",
      "scene": "你开始计算 3.6/4.0 和 3.7/4.0 的命运差异。",
      "description": "进入简历关键词密度。",
      "effects": {},
      "unlocks": ["resume_keyword_density"]
    }
  ],
  "ended": false,
  "ending": null
}
```

### GET `/api/nodes`

返回所有节点配置，用于前端 pipeline 和调试。

## 状态持久化

没有服务端持久化。

没有浏览器持久化。

当前局状态只存在于页面 JavaScript 内存中。每次行动时前端把完整 `state` 发给后端，后端计算并返回下一状态。

## 验收

```bash
go test ./...
go test -bench=. ./...
go run ./cmd/workflowbench run --seed 42
go run ./cmd/server
```

打开 `http://localhost:8080` 能玩完整一局。
