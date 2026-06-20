# WorkflowBench / Metric Hell

一个讽刺学生评价流水线的互动 benchmark 游戏。

项目核心不是预测人生，也不是评价任何真实学生、学校、地区、职业或家庭背景，而是讽刺“人被无限指标化、无限比较、无限追问”的荒诞感。

## 项目声明

本项目不预测真实人生。

本项目不评价真实学校、地区、职业或学生。

本项目所有角色与节点均为虚构或抽象化表达。

本项目不收集个人信息，不登录，不保存玩家结果，不使用数据库。

请不要把本项目结果当成升学、就业或人生建议。如果你玩完觉得焦虑，请关闭它，然后去喝水、睡觉、散步、找朋友聊天。

## 功能

- Go 标准库 HTTP 服务。
- 原生 HTML/CSS/JS 前端。
- 每关提交一个可比较字段，然后由系统生成讽刺性审计判词。
- 无状态 API：浏览器每次把当前 `state` 和本轮 `submission` 发给后端，后端只计算下一步。
- `data/nodes.json` 配置 benchmark 节点。
- CLI 纯文本模拟。
- 单元测试和 benchmark。
- Vercel 部署配置。

## 本地运行 Web 版本

```bash
go run ./cmd/server
```

打开：

```text
http://localhost:8080
```

服务会读取 `PORT` 环境变量；未设置时默认使用 `8080`。

## 运行 CLI 版本

```bash
go run ./cmd/workflowbench run --seed 42
```

## 运行测试

```bash
go test ./...
```

如果本机 Go cache 没有写权限，可以临时指定：

```bash
GOCACHE=/tmp/go-build-cache go test ./...
```

## 运行 benchmark

```bash
go test -bench=. ./...
```

## 部署到 Vercel

本项目包含 `vercel.json` 和 `api/index.go`，用于 Vercel Go Runtime。

推荐流程：

```bash
vercel
```

或连接 GitHub 仓库后在 Vercel 控制台导入项目。

部署后，同一个 URL 会提供：

- `/` 游戏页面
- `/api/new`
- `/api/action`
- `/api/nodes`

注意：Vercel Go Runtime 的行为可能随平台更新变化。如果遇到 Go Runtime 兼容问题，可以使用 Render/Fly.io 按 `go run ./cmd/server` 的方式部署为单 Go Web Service。

## 如何新增节点

编辑：

```text
data/nodes.json
```

新增节点需要：

- `id`：唯一字符串。
- `title`：展示名称。
- `stage`：阶段，例如 `高中`、`大学`、`实习`、`大厂/AI`。
- `scenario`：玩家进入该 bench 时所处的具体人生场景。
- `measurement`：该指标的讽刺口径、范围或解释，例如 GPA 节点要说明 4.0 / 5.0 满分制只是在这里被粗暴折成比较字段。
- `input`：本关输入规格，包含 `type`（`number` 或 `select`）、`prompt`、`placeholder`、`help`。
- `options`：系统可匹配的分数桶或选择档位。每项需要 `id`、`label`、`verdict`、`proof`、`effects`、`unlocks`；数字桶可额外设置 `min` / `max`。
- `questions`：当前 benchmark 的荒诞追问。
- `unlocks`：完成后解锁的后续节点 ID。
- `text_on_enter` / `text_on_pass` / `text_on_fail`：系统提示文案。
- `effects`：对各项指标的影响。
- `absurdity`：额外荒诞度。

新增后运行：

```bash
go test ./...
```

## 内容边界

本项目讽刺指标，不讽刺人。

可以讽刺：

- “系统把专业翻译成岗位关键词。”
- “第一学历过滤器把清北、华五、C9、985、211、双非和 QS100 折成入口字段。”
- “拿到实习后系统继续识别厂牌 tier。”
- “进入 BigFactory 后系统继续追问岗位、转正率、裁员风险和 AI 替代概率。”

不要写：

- 地域歧视。
- 学校羞辱。
- 职业羞辱。
- 家庭背景羞辱。
- 真实政策攻击。
- 真实群体攻击。
- 宿命论判断。
- 诱导焦虑或自我伤害暗示。

禁止输出类似：

- “你废了”
- “没前途”
- “失败人生”
- “某学校不行”
- “某岗位低级”

## 隐私说明

- 不登录。
- 不收集真实姓名。
- 不收集学校。
- 不收集真实城市。
- 不收集家庭收入。
- 不收集手机号。
- 不上传照片。
- 不保存个人结果到服务器。
- 不使用 `localStorage` 保存结果。

刷新页面后当前局会丢失，这是有意设计。
