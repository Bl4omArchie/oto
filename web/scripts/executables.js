document.addEventListener('DOMContentLoaded', () => {
  const executableList = document.getElementById('executableList');
  const template = document.getElementById('executableTemplate');
  const searchInput = document.getElementById('search');
  const addBtn = document.getElementById('addBtn');
  const addForm = document.getElementById('addForm');
  const cancelBtn = document.getElementById('cancelBtn');
  const executableForm = document.getElementById('executableForm');

  let executables = [];

  // Fetch executables from API
  async function loadExecutables() {
    try {
      const res = await fetch('/executables'); // Adjust API endpoint if needed
      executables = await res.json();
      renderExecutables(executables);
    } catch (err) {
      console.error('Failed to load executables:', err);
    }
  }

  // Render executables into the list
  function renderExecutables(list) {
    executableList.innerHTML = '';
    list.forEach(exe => {
      const clone = template.content.cloneNode(true);
      clone.querySelector('.executable-name').textContent = `${exe.name} - ${exe.version}`;
      clone.querySelector('.executable-path').textContent = exe.path;
      clone.querySelector('.executable-description').textContent = exe.description;
      executableList.appendChild(clone);
    });
  }

  // Search filter
  searchInput.addEventListener('input', () => {
    const query = searchInput.value.toLowerCase();
    const filtered = executables.filter(e =>
      e.name.toLowerCase().includes(query) ||
      e.tag.toLowerCase().includes(query) ||
      e.description.toLowerCase().includes(query)
    );
    renderExecutables(filtered);
  });

  // Show add form
  addBtn.addEventListener('click', () => addForm.classList.remove('hidden'));

  // Cancel add form
  cancelBtn.addEventListener('click', () => {
    addForm.classList.add('hidden');
    executableForm.reset();
  });

  // Submit new executable
  executableForm.addEventListener('submit', async (e) => {
    e.preventDefault();
    const formData = new FormData(executableForm);
    const data = Object.fromEntries(formData.entries());

    try {
      const res = await fetch('/executables', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(data)
      });

      if (!res.ok) throw new Error('Failed to add executable');

      const newExe = await res.json();
      executables.push(newExe);
      renderExecutables(executables);
      addForm.classList.add('hidden');
      executableForm.reset();
    } catch (err) {
      console.error(err);
      alert('Error adding executable');
    }
  });

  // Initial load
  loadExecutables();
});
