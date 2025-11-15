// tailwind init
tailwind.config = {
    darkMode: 'class',
}

const pageContent = document.getElementById('page-content');
const pageTitle = document.getElementById('page-title');
const breadcrumb = document.getElementById('breadcrumb');
const sidebarBtns = document.querySelectorAll('.sidebar-btn');

// Function to load HTML page into #page-content
async function loadPage(pageFile) {
  const response = await fetch(`pages/${pageFile}`);
  if (response.ok) {
    const html = await response.text();
    pageContent.innerHTML = html;

    // Update title and breadcrumb
    const title = pageFile.split('.')[0];
    pageTitle.textContent = title.charAt(0).toUpperCase() + title.slice(1);
    breadcrumb.textContent = title.charAt(0).toUpperCase() + title.slice(1);
  } else {
    pageContent.innerHTML = `<p>Error loading page.</p>`;
  }
}

// Handle sidebar clicks
sidebarBtns.forEach(btn => {
  btn.addEventListener('click', () => {
    const pageFile = btn.dataset.page;

    // Load selected page
    loadPage(pageFile);

    // Highlight active button
    sidebarBtns.forEach(b => b.classList.remove('active'));
    btn.classList.add('active');
  });
});

// Night mode toggle
const nightModeBtn = document.getElementById('night-mode-toggle');
nightModeBtn.addEventListener('click', () => {
    document.body.classList.toggle('dark');
});

// Load default page
loadPage('home.html');
sidebarBtns[0].classList.add('active');
