document.addEventListener('DOMContentLoaded', function () {
  const store = window.FacultyStore;
  const listEl = document.getElementById('committeeList');
  const countText = document.getElementById('committeeCountText');
  const emptyCommittee = document.getElementById('emptyCommittee');

  function renderList() {
    const members = store ? store.getState().committeeMembers : [];
    listEl.innerHTML = '';

    // Show/hide empty state
    if (emptyCommittee) {
      emptyCommittee.classList.toggle('hidden', members.length > 0);
    }

    members.forEach(member => {
      const li = document.createElement('li');
      li.className = 'flex items-start gap-3 px-3 py-3 rounded-xl bg-gray-50 hover:bg-gray-100 transition-colors';

      const initials = member.name
        .split(' ')
        .map(p => p[0])
        .join('')
        .slice(0, 2)
        .toUpperCase();

      li.innerHTML = `
        <div class="w-10 h-10 rounded-full bg-gradient-to-br from-orange-primary to-orange-secondary flex items-center justify-center text-white text-sm font-bold flex-shrink-0">${initials}</div>
        <div class="flex-1">
          <p class="font-semibold text-gray-900">${member.name}</p>
          <p class="text-xs text-orange-600 font-medium">${member.role}</p>
          <p class="text-xs text-gray-500 flex items-center gap-1 mt-1">
            <i class="fas fa-envelope text-[10px]"></i>
            ${member.email}
          </p>
        </div>
      `;

      listEl.appendChild(li);
    });

    const count = members.length;
    countText.textContent = `${count} member${count === 1 ? '' : 's'}`;
  }

  // Mobile nav toggle
  document.getElementById('mobile-menu-btn')?.addEventListener('click', function () {
    const menu = document.getElementById('mobile-menu');
    menu?.classList.toggle('hidden');
  });

  // Dark mode support
  (function () {
    const savedTheme = localStorage.getItem('theme');
    const savedSystemTheme = localStorage.getItem('systemTheme');
    const body = document.body;

    function applySystemTheme() {
      if (window.matchMedia && window.matchMedia('(prefers-color-scheme: dark)').matches) {
        body.classList.add('dark-mode');
      } else {
        body.classList.remove('dark-mode');
      }
    }

    if (savedSystemTheme === 'true') {
      applySystemTheme();
      if (window.matchMedia) {
        window.matchMedia('(prefers-color-scheme: dark)').addEventListener('change', applySystemTheme);
      }
    } else if (savedTheme === 'dark') {
      body.classList.add('dark-mode');
    }
  })();

  renderList();
});
