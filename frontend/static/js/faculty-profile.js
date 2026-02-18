document.addEventListener('DOMContentLoaded', function () {
  const store = window.FacultyStore;

  const nameHeading = document.getElementById('facultyNameHeading');
  const deptHeading = document.getElementById('facultyDeptHeading');
  const avatarInitials = document.getElementById('avatarInitials');

  const nameInput = document.getElementById('facultyName');
  const emailInput = document.getElementById('facultyEmail');
  const deptInput = document.getElementById('facultyDept');
  const canManageInput = document.getElementById('canManageCommittee');
  const form = document.getElementById('profileForm');
  const resetBtn = document.getElementById('resetProfileBtn');

  function initialsFromName(name) {
    return name
      .split(' ')
      .map(p => p[0])
      .join('')
      .slice(0, 2)
      .toUpperCase();
  }

  function loadFromStore() {
    const faculty = store.getState().faculty;
    nameHeading.textContent = faculty.name;
    deptHeading.textContent = faculty.department;
    avatarInitials.textContent = initialsFromName(faculty.name);

    nameInput.value = faculty.name;
    emailInput.value = faculty.email;
    deptInput.value = faculty.department;
    canManageInput.checked = !!faculty.canManageCommittee;
  }

  form.addEventListener('submit', function (e) {
    e.preventDefault();
    store.updateFaculty({
      name: nameInput.value.trim(),
      email: emailInput.value.trim(),
      department: deptInput.value.trim(),
      canManageCommittee: canManageInput.checked
    });
    loadFromStore();
    alert('Profile updated. This information will be used across all faculty pages.');
  });

  resetBtn.addEventListener('click', function () {
    if (!confirm('Reset profile to default values?')) return;
    store.resetState();
    loadFromStore();
  });

  // Mobile nav + dark mode
  document.getElementById('mobile-menu-btn')?.addEventListener('click', function () {
    const menu = document.getElementById('mobile-menu');
    menu.classList.toggle('hidden');
  });

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

  loadFromStore();
});
