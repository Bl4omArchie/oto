const API_BASE = "http://localhost:8080/api";
let cmdParamsData = [];
let allFlags = [];
let currentPage = 1;
const pageSize = 20;

const get = id => document.getElementById(id).value.trim();
const check = id => document.getElementById(id).checked;
let binaryInput;

document.addEventListener("DOMContentLoaded", async () => {
  binaryInput = document.getElementById("exec_binary");
  const binaryDrop = document.getElementById("binary-drop");
  const execId = new URLSearchParams(window.location.search).get("exec_id");

  // Binary drag-and-drop
  if (binaryDrop && binaryInput) {
    binaryDrop.addEventListener("click", () => binaryInput.click());
    binaryDrop.addEventListener("dragover", e => { e.preventDefault(); binaryDrop.classList.add("hover"); });
    binaryDrop.addEventListener("dragleave", e => { e.preventDefault(); binaryDrop.classList.remove("hover"); });
    binaryDrop.addEventListener("drop", e => {
      e.preventDefault();
      binaryDrop.classList.remove("hover");
      const file = e.dataTransfer.files[0];
      if (file) {
        binaryInput.files = e.dataTransfer.files;
        binaryDrop.querySelector("p").textContent = file.name;
      }
    });
  }

  if (execId) loadExecInfo(execId);

  await loadValueTypes();
  await loadExecutablesDropdown();
  await loadCommandDropdowns();
  await loadCommandParams();
});

async function createExecutable() {
  const name = get("exec_name");
  const version = get("exec_version");
  const description = get("exec_description");
  const file = binaryInput.files[0];

  if (!file) return alert("Please provide a binary file.");

  const formData = new FormData();
  formData.append("name", name);
  formData.append("version", version);
  formData.append("description", description);
  formData.append("binary", file);

  const res = await fetch(`${API_BASE}/executables`, { method: "POST", body: formData });
  const text = await res.text();
  try { renderTerminal(JSON.stringify(JSON.parse(text), null, 2)); }
  catch { renderTerminal(text); }

  await loadExecutablesDropdown();
}

function createParameter() {
  post("parameters", {
    flag: get("param_flag"),
    exec_id: get("param_execid_select"),
    requires_root: check("param_root"),
    requires_value: check("param_value"),
    value_type: get("param_valuetype"),
    description: get("param_description"),
  });
  loadCommandParams();
}

function createCommand() {
  const selectedFlags = Array.from(document.querySelectorAll("#cmd_params_box div.selected"))
                             .map(d => d.textContent);

  post("commands", {
    name: get("cmd_name"),
    exec_id: get("cmd_execid"),
    description: get("cmd_description"),
    parameter_flags: selectedFlags,
  });
}

async function post(endpoint, body) {
  const res = await fetch(`${API_BASE}/${endpoint}`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(body),
  });
  const text = await res.text();
  try { renderTerminal(JSON.stringify(JSON.parse(text), null, 2)); }
  catch { renderTerminal(text); }
}

function renderTerminal(content) {
  document.getElementById("output").textContent = content;
}

async function loadExecutablesDropdown() {
  const execs = await getJSON("executables");
  const paramSelect = document.getElementById("param_execid_select");
  const cmdSelect = document.getElementById("cmd_execid");

  [paramSelect, cmdSelect].forEach(s => {
    s.innerHTML = '<option value="">-- Select Executable --</option>';
    execs.forEach(e => s.appendChild(Object.assign(document.createElement("option"), { value: e.exec_id, textContent: `${e.name} (${e.version})` })));
  });
}

async function loadValueTypes() {
  try {
    const types = await getJSON("valuetypes");
    const select = document.getElementById("param_valuetype");
    select.innerHTML = '<option value="">-- Select --</option>';
    types.forEach(t => select.appendChild(Object.assign(document.createElement("option"), { value: t, textContent: t })));
  } catch(e) { console.error(e); }
}

async function loadCommandParams() {
  const params = await getJSON("parameters");
  allFlags = params.map(p => p.flag);
  cmdParamsData = [...allFlags];
  currentPage = 1;
  renderCmdParams();
}

function renderCmdParams() {
  const container = document.getElementById("cmd_params_box");
  container.innerHTML = "";
  const start = (currentPage-1) * pageSize;
  const pageData = cmdParamsData.slice(start, start + pageSize);

  pageData.forEach(flag => {
    const div = document.createElement("div");
    div.textContent = flag;
    div.onclick = () => div.classList.toggle("selected");
    container.appendChild(div);
  });

  const pageCount = Math.ceil(cmdParamsData.length / pageSize);
  document.getElementById("page_info").textContent = `${currentPage}/${pageCount}`;
}

function changePage(delta) {
  const pageCount = Math.ceil(cmdParamsData.length / pageSize);
  currentPage = Math.min(Math.max(currentPage + delta, 1), pageCount);
  renderCmdParams();
}

function filterCmdParams() {
  const term = document.getElementById("cmd_params_search").value.toLowerCase();
  cmdParamsData = allFlags.filter(f => f.toLowerCase().includes(term));
  currentPage = 1;
  renderCmdParams();
}

async function getJSON(endpoint) {
  const res = await fetch(`${API_BASE}/${endpoint}`);
  return res.json();
}

async function loadExecutables() {
  const execs = await getJSON("executables");

  const container = document.getElementById("exec_cards_container");
  container.innerHTML = "";

  for (const exec of execs) {
    const [params, cmds] = await Promise.all([
      getJSON(`parameters?exec_id=${exec.exec_id}`),
      getJSON(`commands?exec_id=${exec.exec_id}`)
    ]);

    const card = document.createElement("section");
    card.className = "card exec-info-card";
    card.innerHTML = `
      <h2>${exec.name}</h2>
      <p>${exec.binary}</p>
      <small>v${exec.version}</small>
      <div class="exec-stats">
        <div class="stat-card params">
          <h4>Parameters</h4>
          <p>${params.length}</p>
        </div>
        <div class="stat-card cmds">
          <h4>Commands</h4>
          <p>${cmds.length}</p>
        </div>
      </div>
      <button onclick="openExecSettings('${exec.exec_id}')">Open Settings</button>
    `;
    container.appendChild(card);
  }
}


async function loadExecInfo(execId) {
  try {
    const exec = await get(`executables/${execId}`);
    const params = await get(`parameters?exec_id=${execId}`);
    const cmds = await get(`commands?exec_id=${execId}`);

    document.getElementById("exec_name_display").textContent = exec.name;
    document.getElementById("exec_bin_display").textContent = exec.binary;
    document.getElementById("exec_version_display").textContent = exec.version;

    document.getElementById("exec_args_count").textContent = params.length;
    document.getElementById("exec_cmd_count").textContent = cmds.length;
  } catch (err) {
    console.error("Failed to load exec info:", err);
  }
}

function openExecSettings() {
  window.location.href = "exec_settings.html";
}
