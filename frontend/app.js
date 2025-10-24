// -------------------------
// Executables Page
// -------------------------
async function loadExecutables() {
  const execContainer = document.getElementById('exec-container');
  if (!execContainer) return;

  execContainer.innerHTML = '<p class="col-span-full text-gray-500 text-center">Loading...</p>';

  try {
    const res = await fetch('http://localhost:8080/execs');
    if (!res.ok) throw new Error('Network error');
    const data = await res.json();
    execContainer.innerHTML = '';

    if (!data.length) {
      execContainer.innerHTML = '<p class="col-span-full text-gray-500 text-center">No executables found.</p>';
      return;
    }

    data.forEach(exec => {
      const card = document.createElement('div');
      card.className = 'bg-white rounded-xl shadow p-6 flex flex-col hover:shadow-2xl transition';
      card.innerHTML = `
        <h3 class="text-lg font-semibold mb-1 text-left">${exec.name}</h3>
        <p class="text-sm text-gray-600 mb-1 text-left">Version: ${exec.version}</p>
        <p class="text-gray-700 text-sm mb-1 text-left">Description: ${exec.description || '-'}</p>
        <p class="text-gray-500 text-xs mb-4 text-left">Binary: ${exec.binary}</p>
        <a href="exec-details.html?execId=${encodeURIComponent(exec.exec_id)}&name=${encodeURIComponent(exec.name)}" class="px-4 py-2 bg-gray-200 text-gray-800 rounded hover:bg-gray-300 transition mb-2">View Commands</a>
        <a href="parameters.html?execId=${encodeURIComponent(exec.exec_id)}&name=${encodeURIComponent(exec.name)}" class="px-4 py-2 bg-gray-200 text-gray-800 rounded hover:bg-gray-300 transition">View Parameters</a>
      `;
      execContainer.appendChild(card);
    });

  } catch (err) {
    execContainer.innerHTML = `<p class="col-span-full text-red-500 text-center">${err.message}</p>`;
  }
}

// -------------------------
// Commands Page
// -------------------------
async function loadCommands() {
  const execId = new URLSearchParams(window.location.search).get('execId');
  const execName = new URLSearchParams(window.location.search).get('name');

  const detailsTitle = document.getElementById('details-title');
  const breadcrumbCurrent = document.getElementById('current-exec');
  const commandsContainer = document.getElementById('commands-container');

  if (!commandsContainer) return;

  detailsTitle.textContent = execName;
  breadcrumbCurrent.textContent = execName;
  commandsContainer.innerHTML = '<p class="text-gray-500 text-center">Loading...</p>';

  try {
    const res = await fetch(`http://localhost:8080/cmds/${execId}`);
    if (!res.ok) throw new Error('Network error');
    const cmds = await res.json();
    commandsContainer.innerHTML = '';

    if (!cmds.length) {
      commandsContainer.innerHTML = '<p class="text-gray-500 text-center">No commands found.</p>';
      return;
    }

    cmds.forEach(cmd => {
      const card = document.createElement('div');
      card.className = 'bg-white rounded-xl shadow p-4 hover:shadow-lg transition';
      const paramFlags = Array.isArray(cmd.ParameterFlags) ? cmd.ParameterFlags.join(', ') : cmd.ParameterFlags || 'None';
      card.innerHTML = `
        <h4 class="font-semibold mb-1">${cmd.Name}</h4>
        <p class="text-sm text-gray-600 mb-1">Requires Root: ${cmd.RequiresRoot}</p>
        <p class="text-gray-700 text-sm mb-1">${cmd.Description || '-'}</p>
        <p class="text-gray-500 text-xs">Parameter Flags: ${paramFlags}</p>
      `;
      commandsContainer.appendChild(card);
    });

  } catch (err) {
    commandsContainer.innerHTML = `<p class="text-red-500 text-center">${err.message}</p>`;
  }
}

// -------------------------
// Parameters Page
// -------------------------
async function loadParameters() {
  const execId = new URLSearchParams(window.location.search).get('execId');
  const execName = new URLSearchParams(window.location.search).get('name');

  const paramsTitle = document.getElementById('parameters-title');
  const breadcrumbCurrent = document.getElementById('current-exec-param');
  const paramsContainer = document.getElementById('parameters-container');

  if (!paramsContainer) return;

  paramsTitle.textContent = execName + " - Parameters";
  breadcrumbCurrent.textContent = execName;
  paramsContainer.innerHTML = '<p class="text-gray-500 text-center">Loading...</p>';

  try {
    const res = await fetch(`http://localhost:8080/params/execId=${execId}`);
    if (!res.ok) throw new Error('Network error');
    const params = await res.json();
    paramsContainer.innerHTML = '';

    if (!params.length) {
      paramsContainer.innerHTML = '<p class="text-gray-500 text-center">No parameters found.</p>';
      return;
    }

    params.forEach(param => {
      const card = document.createElement('div');
      card.className = 'bg-white rounded-xl shadow p-4 hover:shadow-lg transition';
      card.innerHTML = `
        <h4 class="font-semibold mb-1">${param.flag}</h4>
        <p class="text-sm text-gray-600 mb-1">Requires Root: ${param.requires_value}</p>
        <p class="text-sm text-gray-600 mb-1">Requires Value: ${param.requires_root}</p>
        <p class="text-gray-700 text-sm mb-1">Value Type: ${param.value_type || '-'}</p>
        <p class="text-sm text-gray-600 mb-1">Depends On: ${param.depends_on || 'None'}</p>
        <p class="text-sm text-gray-600 mb-1">Conflicts With: ${param.conflict_with || 'None'}</p>
        <p class="text-sm text-gray-600 mb-1">Description: ${param.description || '-'}</p>
      `;
      paramsContainer.appendChild(card);
    });

  } catch (err) {
    paramsContainer.innerHTML = `<p class="text-red-500 text-center">${err.message}</p>`;
  }
}

// -------------------------
// Initialize
// -------------------------
window.addEventListener('DOMContentLoaded', () => {
  loadExecutables();
  loadCommands();
  loadParameters();
});
