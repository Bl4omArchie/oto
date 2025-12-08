const searchInput = document.getElementById('search');
const addBtn = document.getElementById('addBtn');
const addForm = document.getElementById('addForm');
const cancelBtn = document.getElementById('cancelBtn');
const executableForm = document.getElementById('executableForm');
const executableList = document.getElementById('executableList');

let executables = [];

// Fetch executables from backend
async function loadExecutables() {
  const res = await fetch('/executables');
  executables = await res.json();
  renderList();
}

// Render filtered list
function renderList() {
  const query = searchInput.value.toLowerCase();
  executableList.innerHTML = '';
  executables
    .filter(b => b.name.toLowerCase().includes(query) || b.tag.toLowerCase().includes(query))
    .forEach(b => {
      const li = document.createElement('li');
      li.className = 'bg-white p-2 rounded shadow flex justify-between';
      li.innerHTML = `<div>
                        <p class="font-bold">${b.name} (${b.tag})</p>
                        <p class="text-sm text-gray-600">${b.version} — ${b.path}</p>
                        <p class="text-gray-700">${b.description}</p>
                      </div>`;
      executableList.appendChild(li);
    });
}

// Show/hide form
addBtn.addEventListener('click', () => addForm.classList.remove('hidden'));
cancelBtn.addEventListener('click', () => addForm.classList.add('hidden'));

// Handle form submission
executableForm.addEventListener('submit', async e => {
  e.preventDefault();
  const formData = new FormData(executableForm);
  const data = Object.fromEntries(formData.entries());
  const res = await fetch('/executables', {
    method: 'POST',
    headers: {'Content-Type': 'application/json'},
    body: JSON.stringify(data)
  });
  if (res.ok) {
    const newExecutable = await res.json();
    executables.push(newExecutable);
    renderList();
    addForm.classList.add('hidden');
    executableForm.reset();
  } else {
    const err = await res.json();
    alert('Error: ' + (err.error || 'Failed to add executable'));
  }
});

// Search input event
searchInput.addEventListener('input', renderList);

// Initial load
loadExecutables();