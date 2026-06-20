let currentResult = null;
let allNodes = [];

const screens = {
  home: document.querySelector("#home"),
  game: document.querySelector("#game"),
  result: document.querySelector("#result"),
};

const metricDefs = [
  ["bench_score", "BenchScore"],
  ["anxiety", "Anxiety"],
  ["selfhood", "Selfhood"],
  ["energy", "Energy"],
  ["curiosity", "Curiosity"],
  ["parent_pressure", "ParentPressure"],
  ["peer_comparison", "PeerComparison"],
  ["escape_index", "EscapeIndex"],
  ["absurdity", "Absurdity"],
];

document.querySelector("#startBtn").addEventListener("click", startGame);
document.querySelector("#restartBtn").addEventListener("click", startGame);
document.querySelector("#actionForm").addEventListener("submit", submitAction);
document.querySelector("#actionSelect").addEventListener("change", updateActionDescription);

async function startGame() {
  allNodes = await fetchJSON("/api/nodes");
  currentResult = await fetchJSON("/api/new", { method: "POST" });
  renderGame();
  showScreen("game");
}

async function submitAction(event) {
  event.preventDefault();
  if (!currentResult || currentResult.ended) return;
  const action = document.querySelector("#actionSelect").value;
  currentResult = await fetchJSON("/api/action", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({
      state: currentResult.state,
      action,
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

  const questions = document.querySelector("#questions");
  questions.innerHTML = "";
  for (const question of node?.questions || []) {
    const li = document.createElement("li");
    li.textContent = question;
    questions.appendChild(li);
  }

  renderActions(actions || []);
  renderMetrics("#metrics", state);
  renderPipeline(state);
  renderEventLog(state);
}

function renderActions(actions) {
  const select = document.querySelector("#actionSelect");
  select.innerHTML = "";
  for (const action of actions) {
    const option = document.createElement("option");
    option.value = action.id;
    option.textContent = action.label;
    option.dataset.description = action.description;
    select.appendChild(option);
  }
  updateActionDescription();
}

function updateActionDescription() {
  const select = document.querySelector("#actionSelect");
  const option = select.options[select.selectedIndex];
  document.querySelector("#actionDescription").textContent = option?.dataset.description || "";
}

function renderMetrics(selector, state) {
  const container = document.querySelector(selector);
  container.innerHTML = "";
  for (const [key, label] of metricDefs) {
    const value = Number(state[key] || 0);
    const item = document.createElement("div");
    item.className = "metric";
    item.innerHTML = `
      <div class="metric-head"><span>${label}</span><strong>${value}</strong></div>
      <div class="bar"><span style="width: ${value}%"></span></div>
    `;
    container.appendChild(item);
  }
}

function renderPipeline(state) {
  const container = document.querySelector("#pipeline");
  container.innerHTML = "";
  for (const node of allNodes) {
    const item = document.createElement("div");
    const completed = state.completed_nodes?.includes(node.id);
    const unlocked = state.unlocked_nodes?.includes(node.id);
    item.className = `pipeline-item ${completed ? "done" : unlocked ? "open" : "locked"}`;
    item.innerHTML = `
      <span>${completed ? "✓" : unlocked ? "⏳" : "🔒"}</span>
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
