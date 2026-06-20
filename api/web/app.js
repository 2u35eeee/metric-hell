let currentResult = null;
let allNodes = [];
let pendingAudit = null;

const screens = {
  home: document.querySelector("#home"),
  game: document.querySelector("#game"),
  result: document.querySelector("#result"),
};

const metricDefs = [
  ["bench_score", "被排序分", "系统把你塞进榜单的顺滑程度。"],
  ["anxiety", "焦虑负载", "为了满足下一张表而持续后台运行的压力。"],
  ["selfhood", "自我保留量", "没有被排名、厂牌和关键词吃掉的那部分自己。"],
  ["energy", "能量余额", "还能不能像人一样睡觉、发呆、恢复。"],
  ["curiosity", "好奇心", "还会不会问“我想知道什么”。"],
  ["parent_pressure", "外部催促压", "来自亲友、默认路径和稳定叙事的合力。"],
  ["peer_comparison", "同辈比较浓度", "把别人的进度条误读成自己的判决书。"],
  ["escape_index", "逃逸指数", "拒绝被单一字段解释的能力。"],
  ["absurdity", "荒诞浓度", "系统越认真，事情越不像人话。"],
];

document.querySelector("#startBtn").addEventListener("click", startGame);
document.querySelector("#restartBtn").addEventListener("click", startGame);
document.querySelector("#submissionForm").addEventListener("submit", submitField);
document.querySelector("#continueBtn").addEventListener("click", continueBenchmark);

async function startGame() {
  pendingAudit = null;
  allNodes = await fetchJSON("/api/nodes");
  currentResult = await fetchJSON("/api/new", { method: "POST" });
  renderGame();
  showScreen("game");
}

async function submitField(event) {
  event.preventDefault();
  if (!currentResult || currentResult.ended || pendingAudit) return;
  const node = currentResult.current_node;
  const submission = buildSubmission(node);
  if (!submission) return;

  currentResult = await fetchJSON("/api/action", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({
      state: currentResult.state,
      submission,
    }),
  });
  pendingAudit = currentResult.audit_record || null;
  if (currentResult.ended) {
    renderResult();
    showScreen("result");
    return;
  }
  renderGame();
}

function continueBenchmark() {
  pendingAudit = null;
  renderGame();
}

async function fetchJSON(url, options = {}) {
  const response = await fetch(url, options);
  const data = await response.json();
  if (!response.ok) {
    throw new Error(data.error || "请求失败");
  }
  return data;
}

function showScreen(name) {
  for (const [screenName, node] of Object.entries(screens)) {
    node.classList.toggle("hidden", screenName !== name);
  }
}

function renderGame() {
  const { state, current_node: node } = currentResult;
  document.querySelector("#studentId").textContent = state.virtual_student_id || "虚构学生";
  document.querySelector("#completionValue").textContent = `${completionPercent(state)}%`;
  renderMetrics("#metrics", state, { compact: true });
  renderPipeline(state, pendingAudit ? null : node?.id);
  renderEventLog(state);

  if (pendingAudit) {
    renderAudit(pendingAudit, node);
    return;
  }
  renderInputNode(node);
}

function renderInputNode(node) {
  showInputMode();
  document.querySelector("#nodeStage").textContent = node?.stage || "系统";
  document.querySelector("#nodeTitle").textContent = node?.title || "系统正在生成下一张表";
  document.querySelector("#nodeText").textContent = node?.scenario || node?.text_on_enter || "";
  document.querySelector("#inputPrompt").textContent = node?.input?.prompt || "请提交一个可比较字段";
  document.querySelector("#inputHelp").textContent = node?.input?.help || "";

  const numberInput = document.querySelector("#numericInput");
  const selectInput = document.querySelector("#optionInput");
  numberInput.classList.toggle("hidden", node?.input?.type !== "number");
  selectInput.classList.toggle("hidden", node?.input?.type !== "select");
  numberInput.value = "";
  numberInput.placeholder = node?.input?.placeholder || "";
  selectInput.innerHTML = "";
  for (const option of node?.options || []) {
    const item = document.createElement("option");
    item.value = option.id;
    item.textContent = option.label;
    selectInput.appendChild(item);
  }

  const hints = document.querySelector("#questionHints");
  hints.innerHTML = "";
  for (const question of (node?.questions || []).slice(0, 3)) {
    const li = document.createElement("li");
    li.textContent = question;
    hints.appendChild(li);
  }
}

function renderAudit(audit, nextNode) {
  showAuditMode();
  document.querySelector("#auditStage").textContent = audit.stage;
  document.querySelector("#auditTitle").textContent = audit.node_title;
  document.querySelector("#auditSubmitted").textContent = audit.submitted_label;
  document.querySelector("#auditVerdict").textContent = audit.verdict;
  document.querySelector("#auditProof").textContent = audit.proof;
  document.querySelector("#nextNodeName").textContent = nextNode?.title || "最终分析报告";

  const changes = document.querySelector("#auditEffects");
  changes.innerHTML = "";
  for (const item of effectPills(audit.effects)) {
    changes.appendChild(item);
  }
}

function showInputMode() {
  document.querySelector("#inputPanel").classList.remove("hidden");
  document.querySelector("#auditPanel").classList.add("hidden");
}

function showAuditMode() {
  document.querySelector("#inputPanel").classList.add("hidden");
  document.querySelector("#auditPanel").classList.remove("hidden");
}

function buildSubmission(node) {
  if (!node) return null;
  if (node.input?.type === "number") {
    const raw = document.querySelector("#numericInput").value.trim();
    const value = Number(raw);
    if (!raw || Number.isNaN(value)) {
      document.querySelector("#inputHelp").textContent = "系统无法解析该字段，请输入一个数字。";
      return null;
    }
    return { node_id: node.id, numeric_value: value };
  }
  return {
    node_id: node.id,
    option_id: document.querySelector("#optionInput").value,
  };
}

function renderMetrics(selector, state, options = {}) {
  const container = document.querySelector(selector);
  container.innerHTML = "";
  for (const [key, label, description] of metricDefs) {
    const value = Number(state[key] || 0);
    const item = document.createElement("div");
    item.className = options.compact ? "metric compact" : "metric";
    const head = document.createElement("div");
    head.className = "metric-head";
    head.innerHTML = `<span>${label}</span><strong>${value}/100</strong>`;
    const bar = document.createElement("div");
    bar.className = "bar";
    const fill = document.createElement("span");
    fill.style.width = `${value}%`;
    bar.appendChild(fill);
    item.append(head);
    if (!options.compact) {
      const copy = document.createElement("p");
      copy.className = "metric-copy";
      copy.textContent = description;
      item.appendChild(copy);
    }
    item.appendChild(bar);
    container.appendChild(item);
  }
}

function renderPipeline(state, currentNodeId) {
  const container = document.querySelector("#pipeline");
  container.innerHTML = "";
  for (const node of allNodes) {
    const completed = state.completed_nodes?.includes(node.id);
    const unlocked = state.unlocked_nodes?.includes(node.id);
    const current = currentNodeId === node.id;
    const item = document.createElement("div");
    item.className = `pipeline-item ${current ? "current" : completed ? "done" : unlocked ? "open" : "locked"}`;
    item.innerHTML = `
      <span>${current ? "→" : completed ? "✓" : unlocked ? "⏳" : "🔒"}</span>
      <div>
        <strong>${node.stage}</strong>
        <p>${node.title}</p>
      </div>
    `;
    container.appendChild(item);
  }
}

function renderEventLog(state) {
  const container = document.querySelector("#eventLog");
  container.innerHTML = "";
  const logs = [...(state.event_log || [])].slice(-6).reverse();
  for (const entry of logs) {
    const item = document.createElement("p");
    item.textContent = entry;
    container.appendChild(item);
  }
}

function renderResult() {
  const { state, ending, audit_record: audit } = currentResult;
  if (audit) {
    pendingAudit = audit;
  }
  document.querySelector("#endingTitle").textContent = "分析结果";
  document.querySelector("#endingType").textContent = ending?.type || "系统分析报告";
  document.querySelector("#finalVerdict").textContent = ending?.system_evaluation || "系统已经生成足够多的字段。";
  document.querySelector("#hiddenEvaluation").textContent = ending?.hidden_evaluation || "这份报告仍然解释不了完整的人。";
  renderConclusions(state);
  renderAuditTrail(state.audit_trail || []);
  renderMetrics("#resultMetrics", state);
}

function renderConclusions(state) {
  const list = document.querySelector("#conclusions");
  list.innerHTML = "";
  const pressure = topMetric(state, ["anxiety", "parent_pressure", "peer_comparison"]);
  const escape = topMetric(state, ["selfhood", "curiosity", "escape_index"]);
  const absurd = Number(state.absurdity || 0);
  const items = [
    `最高压力来源：${metricLabel(pressure.key)} ${pressure.value}/100。系统最擅长把这个字段解释成“还要继续”。`,
    `最荒谬审计：荒诞浓度 ${absurd}/100。系统越认真，越像在用表格表演玄学。`,
    `仍未被解释的部分：${metricLabel(escape.key)} ${escape.value}/100。这里残留了一点没有被字段吃掉的人。`,
  ];
  for (const text of items) {
    const li = document.createElement("li");
    li.textContent = text;
    list.appendChild(li);
  }
}

function renderAuditTrail(records) {
  const list = document.querySelector("#auditTrail");
  list.innerHTML = "";
  for (const record of records) {
    const item = document.createElement("article");
    item.className = "trail-item";
    item.innerHTML = `
      <div class="trail-head">
        <span>${record.turn}. ${record.node_title}</span>
        <strong>${record.submitted_label}</strong>
      </div>
      <p></p>
      <small></small>
    `;
    item.querySelector("p").textContent = record.verdict;
    item.querySelector("small").textContent = record.proof;
    list.appendChild(item);
  }
}

function completionPercent(state) {
  if (!allNodes.length) return "0.0";
  const value = Math.min((state.completed_nodes?.length || 0) / allNodes.length * 100, 99.9);
  return value.toFixed(1);
}

function effectPills(effects = {}) {
  return metricDefs.flatMap(([key, label]) => {
    const value = Number(effects[key] || 0);
    if (!value) return [];
    const pill = document.createElement("span");
    pill.className = value > 0 ? "effect-pill up" : "effect-pill down";
    pill.textContent = `${label} ${value > 0 ? "+" : ""}${value}`;
    return [pill];
  });
}

function topMetric(state, keys) {
  return keys
    .map((key) => ({ key, value: Number(state[key] || 0) }))
    .sort((a, b) => b.value - a.value)[0];
}

function metricLabel(key) {
  return metricDefs.find(([metricKey]) => metricKey === key)?.[1] || key;
}
