document.addEventListener('DOMContentLoaded', async function () {
  const summaryBadges = document.getElementById('progressSummaryBadges');
  const listEl = document.getElementById('acceptedIdeasList');
  const emptyEl = document.getElementById('noAcceptedIdeas');

  const token = localStorage.getItem('authToken');
  if (!token) {
    window.location.href = '/login.html';
    return;
  }

  /* ---------------- Dark mode ---------------- */
  (function applyDarkMode() {
    const savedTheme = localStorage.getItem('theme');
    const savedSystemTheme = localStorage.getItem('systemTheme');
    const body = document.body;

    function applySystemTheme() {
      body.classList.toggle(
        'dark-mode',
        window.matchMedia('(prefers-color-scheme: dark)').matches
      );
    }

    if (savedSystemTheme === 'true') {
      applySystemTheme();
      window.matchMedia('(prefers-color-scheme: dark)')
        .addEventListener('change', applySystemTheme);
    } else if (savedTheme === 'dark') {
      body.classList.add('dark-mode');
    }
  })();

  /* ---------------- Fetch progress ---------------- */
  async function loadProgress() {
    const res = await fetch('/api/faculty/progress', {
      headers: { Authorization: 'Bearer ' + token }
    });

    if (!res.ok) {
      alert('Failed to load incubation progress');
      return [];
    }

    return res.json();
  }

  /* ---------------- Render ---------------- */
  function render(ideas) {
    if (!ideas || ideas.length === 0) {
      emptyEl.classList.remove('hidden');
      listEl.innerHTML = '';
      summaryBadges.innerHTML = '';
      return;
    }

    emptyEl.classList.add('hidden');

    /* ---- Summary badges ---- */
    const totalProgress = ideas.reduce(
      (sum, i) => sum + (i.progress || 0),
      0
    );

    const avgProgress = Math.round(totalProgress / ideas.length);
    const domains = [...new Set(ideas.map(i => i.domain || 'Unspecified'))];

    summaryBadges.innerHTML = `
      <span class="inline-flex items-center gap-2 px-3 py-1 rounded-full bg-gray-900 text-white">
        <i class="fas fa-lightbulb"></i>
        ${ideas.length} accepted idea${ideas.length === 1 ? '' : 's'}
      </span>

      <span class="inline-flex items-center gap-2 px-3 py-1 rounded-full bg-orange-50 text-orange-800 border border-orange-200">
        <i class="fas fa-percent"></i>
        Average progress: ${avgProgress}%
      </span>

      <span class="inline-flex items-center gap-2 px-3 py-1 rounded-full bg-blue-50 text-blue-800 border border-blue-200">
        <i class="fas fa-tag"></i>
        ${domains.length} domain${domains.length === 1 ? '' : 's'}
      </span>
    `;

    /* ---- Cards ---- */
    listEl.innerHTML = '';

    ideas.forEach(idea => {
      const progress = Math.min(100, Math.max(0, idea.progress || 0));
      const stage = idea.stage || 'ideation';
      const incubationStarted = progress > 0;

      const card = document.createElement('div');
      card.className = 'glass-card rounded-2xl p-5 bg-white/80 border border-gray-200/70';

      card.innerHTML = `
        <div class="flex flex-col md:flex-row md:items-center md:justify-between gap-3 mb-3">
          <div>
            <h2 class="text-lg font-semibold text-gray-900 mb-1">${idea.title}</h2>
            <p class="text-xs text-gray-600">${idea.student}</p>
            <p class="text-[11px] text-gray-500 mt-1 flex items-center gap-1">
              <i class="far fa-calendar text-[10px]"></i>
              Accepted on ${new Date(idea.accepted_at).toLocaleDateString()}
            </p>
          </div>

          <div class="text-right text-xs text-gray-600">
            <p class="font-semibold text-gray-800 mb-1 flex items-center justify-end gap-1">
              <i class="fas fa-layer-group text-orange-primary"></i>
              ${stage}
            </p>

            ${
              incubationStarted
                ? `<p><span class="font-semibold">${progress}%</span> complete</p>`
                : `<p class="italic text-gray-500">Incubation not started</p>`
            }
          </div>
        </div>

        <div class="mb-2">
          <div class="progress-bar">
            <div class="progress-fill" style="width:${progress}%"></div>
          </div>
        </div>

        <div class="flex flex-wrap items-center justify-between gap-2 text-[11px] text-gray-600 mt-1">
          <span class="inline-flex items-center gap-1 px-2 py-0.5 rounded-full bg-gray-100 text-gray-700">
            <i class="fas fa-tag text-[10px]"></i>${idea.domain || 'Unspecified domain'}
          </span>

          <button
            class="inline-flex items-center gap-1 px-3 py-1 rounded-full border border-gray-300 hover:border-orange-primary hover:text-orange-primary text-[11px]"
            data-open-idea="${idea.submission_id}">
            <i class="fas fa-eye"></i>
            View idea details
          </button>
        </div>
      `;

      listEl.appendChild(card);
    });
  }

  /* ---------------- Navigation ---------------- */
  listEl.addEventListener('click', function (e) {
    const btn = e.target.closest('[data-open-idea]');
    if (!btn) return;

    const submissionId = btn.dataset.openIdea;
    if (!submissionId) {
      console.error('Missing submission_id on button');
      return;
    }

    window.location.href =
      `/faculty/faculty-idea.html?id=${encodeURIComponent(submissionId)}`;
  });

  document.getElementById('mobile-menu-btn')
    ?.addEventListener('click', () =>
      document.getElementById('mobile-menu')?.classList.toggle('hidden')
    );

  
  const ideas = await loadProgress();
  render(ideas);
});