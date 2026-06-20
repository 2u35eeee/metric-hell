# WorkflowBench / Metric Hell MVP Spec

## 目标

实现一个可上线试玩的 Go 前后端小游戏。玩家打开 URL 后，每关提交一个可比较字段，系统根据字段生成讽刺性审计判词，体验 benchmark 不断增殖的荒诞感。

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

每个 benchmark 节点是一张要求玩家提交字段的审计表：

- `scenario` 描述玩家当前处在什么具体场景里。
- `measurement` 解释该指标到底在讽刺什么、范围/口径是什么。
- `input` 描述本关输入，`type` 为 `number` 或 `select`，并提供 `prompt` / `placeholder` / `help`。
- `options` 描述系统可匹配的分数桶或档位，每项包含 `label`、`verdict`、`proof`、`effects`、`unlocks`；数字桶可设置 `min` / `max`。
- `questions` 是系统追问，用来制造“评价继续增殖”的荒诞感。
- 高考节点使用数字输入和分数桶；GPA、学校层次、厂牌、岗位、风险等节点使用档位选择。

前端必须先让玩家提交字段，再展示“系统判词 / 证明材料 / 指标变化 / 下一张表”。
指标面板是辅助信息，最终结果页必须展示分析报告和路径回放。

## API

### POST `/api/new`

创建一局虚构学生状态。

### POST `/api/action`

请求：

```json
{
  "state": {},
  "submission": {
    "node_id": "gaokao_score",
    "numeric_value": 701
  }
}
```

响应：

```json
{
  "state": {},
  "audit_record": {
    "node_id": "gaokao_score",
    "node_title": "高考成绩 Benchmark",
    "submitted_label": "701",
    "verdict": "你要是 700 以上还认识这个开发者？系统怀疑样本来源异常。",
    "proof": "系统把你标记为罕见样本，同时继续要求更多可比较字段。",
    "effects": {},
    "unlocks": ["province_rank"]
  },
  "current_node": {},
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
