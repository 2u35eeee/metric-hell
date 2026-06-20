let currentResult = null;
let allNodes = [];
let selectedAction = null;

const screens = {
  home: document.querySelector("#home"),
  game: document.querySelector("#game"),
  result: document.querySelector("#result"),
};

const metricDefs = [
  ["bench_score", "被排序分", "系统把你塞进榜单的顺滑程度；越高越像一行合格数据。"],
  ["anxiety", "焦虑负载", "为了满足下一张表而持续后台运行的压力。"],
  ["selfhood", "自我保留量", "没有被排名、厂牌和关键词吃掉的那部分自己。"],
  ["energy", "能量余额", "还能不能像人一样睡觉、发呆、恢复。"],
  ["curiosity", "好奇心", "还会不会问“我想知道什么”，而不只是“别人要什么”。"],
  ["parent_pressure", "外部催促压", "来自亲友、默认路径和稳定叙事的合力。"],
  ["peer_comparison", "同辈比较浓度", "越高越容易把别人的进度条误读成自己的判决书。"],
  ["escape_index", "逃逸指数", "越高越能拒绝被单一字段解释。"],
  ["absurdity", "荒诞浓度", "系统越认真，事情越不像人话。"],
];

document.querySelector("#startBtn").addEventListener("click", startGame);
document.querySelector("#restartBtn").addEventListener("click", startGame);
document.querySelector("#actionForm").addEventListener("submit", submitAction);

async function startGame() {
  allNodes = await fetchJSON("/api/nodes");
  currentResult = await fetchJSON("/api/new", { method: "POST" });
  renderGame();
  showScreen("game");
}

async function submitAction(event) {
  event.preventDefault();
  if (!currentResult || currentResult.ended) return;
  if (!selectedAction) return;
  currentResult = await fetchJSON("/api/action", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({
      state: currentResult.state,
      action: selectedAction,
    }),
  });
  if (currentResult.ended) {
    renderResult();
    showScreen("result");
    return;
  }
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
  const { state, current_node: node, actions } = currentResult;
  document.querySelector("#studentId").textContent = state.virtual_student_id || "虚构学生";
  document.querySelector("#completionValue").textContent = `${completionPercent(state)}%`;
  document.querySelector("#nodeStage").textContent = node?.stage || state.stage;
  document.querySelector("#nodeTitle").textContent = node?.title || "系统正在生成下一张表";
  document.querySelector("#nodeText").textContent = node?.text_on_enter || "提示：人生完成度无法达到 100%，因为系统已发现新的评价维度。";
  setInfoBlock("#nodeScenario", "当前场景", node?.scenario);
  setInfoBlock("#nodeMeasurement", "指标口径", node?.measurement);

  const questions = document.querySelector("#questions");
  questions.innerHTML = "";
  for (const question of node?.questions || []) {
    const li = document.createElement("li");
    li.textContent = question;
    questions.appendChild(li);
  }

  renderActions(actions || []);
  renderMetrics("#metrics", state);
  renderPipeline(state, node?.id);
  renderEventLog(state);
}

function renderActions(actions) {
  const container = document.querySelector("#actions");
  const submit = document.querySelector("#submitAction");
  container.innerHTML = "";
  selectedAction = actions[0]?.id || null;
  submit.disabled = !selectedAction;

  for (const action of actions) {
    const button = document.createElement("button");
    button.type = "button";
    button.className = "action-card";
    button.dataset.actionId = action.id;
    button.setAttribute("aria-pressed", String(action.id === selectedAction));

    const title = document.createElement("strong");
    title.textContent = action.label;

    const scene = document.createElement("p");
    scene.className = "action-scene";
    scene.textContent = action.scene || action.description;

    const description = document.createElement("p");
    description.className = "action-description";
    description.textContent = action.description;

    const meta = document.createElement("div");
    meta.className = "action-meta";
    for (const item of actionEffectPreview(action.effects)) {
      meta.appendChild(item);
    }
    for (const titleText of unlockTitles(action.unlocks)) {
      const pill = document.createElement("span");
      pill.className = "next-pill";
      pill.textContent = `下一站：${titleText}`;
      meta.appendChild(pill);
    }

    button.append(title, scene, description, meta);
    button.addEventListener("click", () => selectAction(action.id));
    container.appendChild(button);
  }

  selectAction(selectedAction);
}

function selectAction(actionId) {
  selectedAction = actionId;
  for (const card of document.querySelectorAll(".action-card")) {
    const selected = card.dataset.actionId === actionId;
    card.classList.toggle("selected", selected);
    card.setAttribute("aria-pressed", String(selected));
  }
}

function renderMetrics(selector, state) {
  const container = document.querySelector(selector);
  container.innerHTML = "";
  for (const [key, label, description] of metricDefs) {
    const value = Number(state[key] || 0);
    const item = document.createElement("div");
    item.className = "metric";
    const head = document.createElement("div");
    head.className = "metric-head";
    head.innerHTML = `<span>${label}</span><strong>${value}/100</strong>`;
    const copy = document.createElement("p");
    copy.className = "metric-copy";
    copy.textContent = description;
    const bar = document.createElement("div");
    bar.className = "bar";
    const fill = document.createElement("span");
    fill.style.width = `${value}%`;
    bar.appendChild(fill);
    item.append(head, copy, bar);
    container.appendChild(item);
  }
}

function renderPipeline(state, currentNodeId) {
  const container = document.querySelector("#pipeline");
  container.innerHTML = "";
  for (const node of allNodes) {
    const item = document.createElement("div");
    const completed = state.completed_nodes?.includes(node.id);
    const unlocked = state.unlocked_nodes?.includes(node.id);
    const current = currentNodeId === node.id;
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
  const logs = [...(state.event_log || [])].slice(-8).reverse();
  for (const entry of logs) {
    const item = document.createElement("p");
    item.textContent = entry;
    container.appendChild(item);
  }
}

function renderResult() {
  const { state, ending } = currentResult;
  document.querySelector("#endingTitle").textContent = `【${ending.title}】`;
  document.querySelector("#endingType").textContent = ending.type;
  document.querySelector("#systemEvaluation").textContent = ending.system_evaluation;
  document.querySelector("#hiddenEvaluation").textContent = ending.hidden_evaluation;
  renderMetrics("#resultMetrics", state);
}

function completionPercent(state) {
  if (!allNodes.length) return "0.0";
  const value = Math.min((state.completed_nodes?.length || 0) / allNodes.length * 100, 99.9);
  return value.toFixed(1);
}

function setInfoBlock(selector, label, value) {
  const node = document.querySelector(selector);
  node.classList.toggle("hidden", !value);
  node.innerHTML = "";
  if (!value) return;
  const strong = document.createElement("strong");
  strong.textContent = label;
  const copy = document.createElement("span");
  copy.textContent = value;
  node.append(strong, copy);
}

function actionEffectPreview(effects = {}) {
  return metricDefs.flatMap(([key, label]) => {
    const value = Number(effects[key] || 0);
    if (!value) return [];
    const pill = document.createElement("span");
    pill.className = value > 0 ? "effect-pill up" : "effect-pill down";
    pill.textContent = `${label} ${value > 0 ? "+" : ""}${value}`;
    return [pill];
  });
}

function unlockTitles(unlocks = []) {
  return unlocks.map((id) => allNodes.find((node) => node.id === id)?.title || id);
}
