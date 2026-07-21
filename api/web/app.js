let currentResult = null;
let allNodes = [];
let latestReceipt = null;
let isSubmitting = false;

const screens = {
  home: document.querySelector("#home"),
  game: document.querySelector("#game"),
  result: document.querySelector("#result"),
};

const stageSteps = [
  { label: "高中", stages: ["高中"] },
  { label: "大学", stages: ["大学", "简历"] },
  { label: "实习", stages: ["实习"] },
  { label: "上班", stages: ["大厂", "大厂/AI"] },
  { label: "系统报错", stages: ["逃逸"] },
];

const gameMetrics = [
  ["bench_score", "系统满意度"],
  ["anxiety", "脑内后台进程"],
  ["selfhood", "还像自己的程度"],
];

const resultMetrics = [
  ["bench_score", "系统满意度"],
  ["anxiety", "焦虑负载"],
  ["selfhood", "自我保留"],
  ["energy", "能量余额"],
  ["escape_index", "逃逸指数"],
  ["absurdity", "荒诞浓度"],
];

document.querySelector("#startBtn").addEventListener("click", startGame);
document.querySelector("#restartBtn").addEventListener("click", startGame);
document.querySelector("#playAgainBtn").addEventListener("click", startGame);
for (const button of document.querySelectorAll("[data-go-home]")) {
  button.addEventListener("click", () => showScreen("home"));
}

async function startGame() {
  setHomeBusy(true);
  clearErrors();
  try {
    [allNodes, currentResult] = await Promise.all([
      fetchJSON("/api/nodes"),
      fetchJSON("/api/new", { method: "POST" }),
    ]);
    latestReceipt = null;
    renderGame();
    showScreen("game");
  } catch (error) {
    document.querySelector("#homeError").textContent = `没启动起来：${error.message}`;
  } finally {
    setHomeBusy(false);
  }
}

async function chooseOption(optionID) {
  if (!currentResult || currentResult.ended || isSubmitting) return;
  isSubmitting = true;
  clearErrors();
  setChoicesBusy(optionID);

  try {
    const node = currentResult.current_node;
    currentResult = await fetchJSON("/api/action", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({
        state: currentResult.state,
        submission: { node_id: node.id, option_id: optionID },
      }),
    });
    latestReceipt = currentResult.audit_record || null;

    if (currentResult.ended) {
      renderResult();
      showScreen("result");
      return;
    }

    renderGame();
    window.scrollTo({ top: 0, behavior: "smooth" });
  } catch (error) {
    document.querySelector("#choiceError").textContent = `系统卡住了：${error.message}`;
    renderChoices(currentResult.current_node);
  } finally {
    isSubmitting = false;
  }
}

async function fetchJSON(url, options = {}) {
  const response = await fetch(url, options);
  const data = await response.json();
  if (!response.ok) throw new Error(data.error || "请求失败");
  return data;
}

function showScreen(name) {
  for (const [screenName, node] of Object.entries(screens)) {
    node.classList.toggle("hidden", screenName !== name);
  }
  window.scrollTo({ top: 0, behavior: "auto" });
}

function renderGame() {
  const { state, current_node: node } = currentResult;
  document.querySelector("#studentId").textContent = state.virtual_student_id || "虚构学生";
  document.querySelector("#turnLabel").textContent = `第 ${state.turn + 1} 轮`;
  document.querySelector("#nodeStage").textContent = node?.stage || "系统";
  document.querySelector("#nodeCount").textContent = `材料 ${pad(state.turn + 1)} / ${allNodes.length}`;
  document.querySelector("#nodeTitle").textContent = cleanTitle(node?.title || "系统还在想问题");
  document.querySelector("#nodeText").textContent = node?.scenario || node?.text_on_enter || "";
  document.querySelector("#choicePrompt").textContent = humanPrompt(node?.input?.prompt || "系统想把这段经历放进哪一格？");
  document.querySelector("#nextQuestion").textContent = node?.questions?.[0]
    ? `它已经在准备下一句：${node.questions[0]}`
    : "";

  renderReceipt();
  renderChoices(node);
  renderStageRail(node?.stage, state);
  renderMetrics("#metrics", state, gameMetrics);
  renderEventLog(state);
}

function renderChoices(node) {
  const list = document.querySelector("#choiceList");
  list.innerHTML = "";
  for (const [index, option] of (node?.options || []).entries()) {
    const button = document.createElement("button");
    button.type = "button";
    button.className = "choice-button";
    button.dataset.optionId = option.id;
    button.innerHTML = `
      <span class="choice-index">${String.fromCharCode(65 + index)}</span>
      <span class="choice-label"></span>
      <span class="choice-arrow" aria-hidden="true">→</span>
    `;
    button.querySelector(".choice-label").textContent = option.label;
    button.addEventListener("click", () => chooseOption(option.id));
    list.appendChild(button);
  }
}

function setChoicesBusy(selectedID) {
  for (const button of document.querySelectorAll(".choice-button")) {
    button.disabled = true;
    if (button.dataset.optionId === selectedID) {
      button.classList.add("selected");
      button.querySelector(".choice-arrow").textContent = "…";
    }
  }
}

function renderReceipt() {
  const box = document.querySelector("#receipt");
  box.classList.toggle("hidden", !latestReceipt);
  if (!latestReceipt) return;
  document.querySelector("#receiptText").textContent = latestReceipt.verdict;
}

function renderStageRail(stage, state) {
  const currentIndex = Math.max(0, stageSteps.findIndex((item) => item.stages.includes(stage)));
  const rail = document.querySelector("#stageRail");
  rail.innerHTML = "";
  for (const [index, step] of stageSteps.entries()) {
    const item = document.createElement("div");
    item.className = index < currentIndex ? "done" : index === currentIndex ? "current" : "";
    item.innerHTML = `<span>${pad(index + 1)}</span><strong>${step.label}</strong>`;
    rail.appendChild(item);
  }
  rail.setAttribute("aria-label", `当前阶段：${stage}，已完成 ${state.completed_nodes?.length || 0} 个 benchmark`);
}

function renderMetrics(selector, state, definitions) {
  const container = document.querySelector(selector);
  container.innerHTML = "";
  for (const [key, label] of definitions) {
    const value = Number(state[key] || 0);
    const item = document.createElement("div");
    item.className = "metric";
    item.innerHTML = `
      <div class="metric-head"><span>${label}</span><strong>${value}</strong></div>
      <div class="metric-track"><span></span></div>
    `;
    item.querySelector(".metric-track span").style.width = `${value}%`;
    container.appendChild(item);
  }
}

function renderEventLog(state) {
  const container = document.querySelector("#eventLog");
  container.innerHTML = "";
  const logs = [...(state.event_log || [])].slice(-3).reverse();
  for (const entry of logs) {
    const item = document.createElement("p");
    item.textContent = shortenLog(entry);
    container.appendChild(item);
  }
}

function renderResult() {
  const { state, ending } = currentResult;
  document.querySelector("#endingTitle").textContent = ending?.title || "系统解释失败";
  document.querySelector("#endingType").textContent = ending?.type || "报告未定义";
  document.querySelector("#finalVerdict").textContent = ending?.system_evaluation || "系统已经生成足够多的字段。";
  document.querySelector("#hiddenEvaluation").textContent = ending?.hidden_evaluation || "这份报告仍然解释不了完整的人。";
  renderMetrics("#resultMetrics", state, resultMetrics);
  renderConclusions(state);
  renderAuditTrail(state.audit_trail || []);
}

function renderConclusions(state) {
  const list = document.querySelector("#conclusions");
  list.innerHTML = "";
  const pressure = topMetric(state, ["anxiety", "parent_pressure", "peer_comparison"]);
  const alive = topMetric(state, ["selfhood", "curiosity", "escape_index"]);
  const items = [
    `${metricName(pressure.key)}最高，${pressure.value}/100。系统把压力当成了继续生成表格的理由。`,
    `荒诞浓度 ${state.absurdity}/100。问题越来越精确，答案却没有更接近一个人。`,
    `${metricName(alive.key)}还剩 ${alive.value}/100。这部分没有被排名、厂牌和岗位名称拿走。`,
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
      <span class="trail-number">${pad(record.turn)}</span>
      <div><strong></strong><p></p><small></small></div>
    `;
    item.querySelector("strong").textContent = `${cleanTitle(record.node_title)}：${record.submitted_label}`;
    item.querySelector("p").textContent = record.verdict;
    item.querySelector("small").textContent = record.proof;
    list.appendChild(item);
  }
}

function humanPrompt(prompt) {
  return prompt
    .replace(/^你的/, "档案里的")
    .replace(/^你/, "这个虚构学生")
    .replace(/被系统/g, "会被系统");
}

function cleanTitle(title) {
  return title.replace(/\s*Benchmark$/i, "").replace(/BigFactory/g, "大厂");
}

function shortenLog(entry) {
  return entry
    .replace(/^系统检测到你完成了/, "做完")
    .replace(/，已生成更细评价指标：/, "，系统又加了：");
}

function topMetric(state, keys) {
  return keys
    .map((key) => ({ key, value: Number(state[key] || 0) }))
    .sort((a, b) => b.value - a.value)[0];
}

function metricName(key) {
  const names = {
    anxiety: "焦虑负载",
    parent_pressure: "外部催促",
    peer_comparison: "同辈比较",
    selfhood: "自我保留",
    curiosity: "好奇心",
    escape_index: "逃逸指数",
  };
  return names[key] || key;
}

function setHomeBusy(busy) {
  const button = document.querySelector("#startBtn");
  button.disabled = busy;
  button.querySelector("span").textContent = busy ? "正在捏造档案…" : "生成一份虚构档案";
}

function clearErrors() {
  document.querySelector("#homeError").textContent = "";
  document.querySelector("#choiceError").textContent = "";
}

function pad(value) {
  return String(value).padStart(2, "0");
}
