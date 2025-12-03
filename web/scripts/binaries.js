const searchInput = document.getElementById('search');
const addBtn = document.getElementById('addBtn');
const addForm = document.getElementById('addForm');
const cancelBtn = document.getElementById('cancelBtn');
const binaryForm = document.getElementById('binaryForm');
const binaryList = document.getElementById('binaryList');

let binaries = [];

// Fetch binaries from backend
async function loadBinaries() {
  const res = await fetch('/binaries');
  binaries = await res.json();
  renderList();
}

// Render filtered list
function renderList() {
  const query = searchInput.value.toLowerCase();
  binaryList.innerHTML = '';
  binaries
    .filter(b => b.name.toLowerCase().includes(query) || b.tag.toLowerCase().includes(query))
    .forEach(b => {
      const li = document.createElement('li');
      li.className = 'bg-white p-2 rounded shadow flex justify-between';
      li.innerHTML = `<div>
                        <p class="font-bold">${b.name} (${b.tag})</p>
                        <p class="text-sm text-gray-600">${b.version} — ${b.path}</p>
                        <p class="text-gray-700">${b.description}</p>
                      </div>`;
      binaryList.appendChild(li);
    });
}

// Show/hide form
addBtn.addEventListener('click', () => addForm.classList.remove('hidden'));
cancelBtn.addEventListener('click', () => addForm.classList.add('hidden'));

// Handle form submission
binaryForm.addEventListener('submit', async e => {
  e.preventDefault();
  const formData = new FormData(binaryForm);
  const data = Object.fromEntries(formData.entries());
  const res = await fetch('/binaries', {
    method: 'POST',
    headers: {'Content-Type': 'application/json'},
    body: JSON.stringify(data)
  });
  if (res.ok) {
    const newBinary = await res.json();
    binaries.push(newBinary);
    renderList();
    addForm.classList.add('hidden');
    binaryForm.reset();
  } else {
    const err = await res.json();
    alert('Error: ' + (err.error || 'Failed to add binary'));
  }
});

// Search input event
searchInput.addEventListener('input', renderList);

// Initial load
loadBinaries();