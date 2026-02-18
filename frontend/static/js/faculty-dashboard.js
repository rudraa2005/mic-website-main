document.addEventListener('DOMContentLoaded', init);

const API = {
  profile: '/api/profile/me',
  ideas: '/api/faculty/reviews',
  progress: '/api/faculty/progress',
  events: '/api/faculty/events/invitations',
  rsvp: id => `/api/faculty/events/invitations/${id}/rsvp`
};

let state = {
  faculty: null,
  ideas: [],
  events: [],
  ideaFilter: 'all'
};

/* -------------------- Helpers -------------------- */

function authHeaders() {
  const token = localStorage.getItem('authToken');
  if (!token) window.location.href = '/login.html';
  return { Authorization: `Bearer ${token}` };
}

async function fetchJSON(url, options = {}) {
  const res = await fetch(url, {
    headers: authHeaders(),
    ...options
  });
  if (!res.ok) throw new Error(await res.text());
  return res.json();
}

/* -------------------- Init -------------------- */

async function init() {
  try {
    await Promise.all([
      loadFaculty(),
      loadIdeas(),
      loadEvents()
    ]);

    renderFaculty();
    renderIdeaStats();
    renderIdeasTable();
    renderEvents();
    wireUI();
  } catch (err) {
    console.error(err);
    alert('Failed to load faculty dashboard');
  }
}

/* -------------------- Loaders -------------------- */

async function loadFaculty() {
  state.faculty = await fetchJSON(API.profile);
}

async function loadIdeas() {
  try {
    const response = await fetchJSON(API.ideas);
    state.ideas = Array.isArray(response) ? response : [];
  } catch (e) {
    console.error('Failed to load ideas:', e);
    state.ideas = [];
  }
}

async function loadEvents() {
  try {
    const response = await fetchJSON(API.events);
    state.events = Array.isArray(response) ? response : [];
  } catch (e) {
    console.error('Failed to load events:', e);
    state.events = [];
  }
}

/* -------------------- Faculty -------------------- */

function renderFaculty() {
  const f = state.faculty;
  document.getElementById('facultyNameText').textContent = f.name;
  document.getElementById('facultyEmailText').textContent = f.email;
  document.getElementById('facultyDepartmentText').textContent = f.department;
}

/* -------------------- Ideas -------------------- */

function renderIdeaStats() {
  // Use admin_approved as the "pending faculty review" status
  const stats = {
    total: state.ideas.length,
    pending: state.ideas.filter(i => i.status === 'admin_approved').length,
    approved: state.ideas.filter(i => i.status === 'approved').length,
    rejected: state.ideas.filter(i => i.status === 'rejected').length
  };

  document.getElementById('statTotalIdeas').textContent = stats.total;
  document.getElementById('statPendingIdeas').textContent = stats.pending;
  document.getElementById('statApprovedIdeas').textContent = stats.approved;
  document.getElementById('statRejectedIdeas').textContent = stats.rejected;

  // Update glance counts
  const ideasToReview = document.getElementById('ideasToReviewCount');
  if (ideasToReview) ideasToReview.textContent = stats.pending;
}

function renderIdeasTable() {
  const tbody = document.getElementById('ideasTableBody');
  tbody.innerHTML = '';

  // Dashboard shows only the 5 most recent ideas (no filtering)
  const ideas = state.ideas.slice(0, 5);

  ideas.forEach(i => {
    const tr = document.createElement('tr');
    tr.dataset.id = i.id;
    tr.className = 'idea-row cursor-pointer border-b';

    tr.innerHTML = `
      <td class="py-3 px-4">
        <div class="font-semibold">${i.title}</div>
        <div class="text-xs text-gray-500">${i.submitted_on}</div>
      </td>
      <td class="py-3 px-4">${i.student}</td>
      <td class="py-3 px-4">${renderIdeaStatus(i.status)}</td>
      <td class="py-3 px-4">
        ${i.requires_review ? badge('Yes', 'orange') : badge('No', 'gray')}
      </td>
    `;

    tr.addEventListener('click', () => {
      window.location.href = `faculty-idea.html?id=${i.id}`;
    });

    tbody.appendChild(tr);
  });
}

function renderIdeaStatus(status) {
  const map = {
    admin_approved: badge('Pending Faculty Review', 'orange'),
    approved: badge('Approved', 'green'),
    rejected: badge('Rejected', 'red'),
    submitted: badge('Pending Admin', 'gray'),
    admin_rejected: badge('Admin Rejected', 'red')
  };
  return map[status] || badge(status, 'gray');
}

/* -------------------- Events -------------------- */

function renderEvents() {
  const container = document.getElementById('eventsList');
  container.innerHTML = '';

  const counts = {
    pending: 0,
    accepted: 0,
    declined: 0
  };

  state.events.forEach(e => counts[e.status]++);

  document.getElementById('acceptedEventsCount').textContent = counts.accepted;
  document.getElementById('pendingEventsCount').textContent = counts.pending;
  document.getElementById('declinedEventsCount').textContent = counts.declined;

  state.events.forEach(event => {
    const card = document.createElement('div');
    card.className = 'event-card rounded-xl border p-4';

    card.innerHTML = `
      <div class="flex justify-between mb-2">
        <h3 class="font-semibold">${event.title}</h3>
        ${badge(event.status, statusColor(event.status))}
      </div>
      <div class="text-xs text-gray-600 mb-2">
        <i class="far fa-calendar"></i> ${event.event_date}
      </div>
      <div class="flex gap-2">
        <button class="accept">Accept</button>
        <button class="decline">Decline</button>
      </div>
    `;

    card.querySelector('.accept').onclick = () =>
      updateRSVP(event.id, 'accepted');

    card.querySelector('.decline').onclick = () =>
      updateRSVP(event.id, 'declined');

    container.appendChild(card);
  });
}

async function updateRSVP(id, status) {
  await fetchJSON(API.rsvp(id), {
    method: 'POST',
    body: JSON.stringify({ status }),
    headers: {
      ...authHeaders(),
      'Content-Type': 'application/json'
    }
  });

  await loadEvents();
  renderEvents();
}

/* -------------------- UI Wiring -------------------- */

function wireUI() {
  document.querySelectorAll('[data-idea-filter]').forEach(btn => {
    btn.onclick = () => {
      state.ideaFilter = btn.dataset.ideaFilter;
      renderIdeasTable();
    };
  });

  document.getElementById('mobile-menu-btn')?.addEventListener('click', () => {
    document.getElementById('mobile-menu')?.classList.toggle('hidden');
  });
}

/* -------------------- UI Utils -------------------- */

function badge(text, color) {
  const map = {
    orange: 'bg-orange-100 text-orange-800',
    green: 'bg-emerald-100 text-emerald-800',
    red: 'bg-rose-100 text-rose-800',
    gray: 'bg-gray-100 text-gray-700'
  };
  return `<span class="px-2 py-1 rounded-full text-xs ${map[color]}">${text}</span>`;
}

function statusColor(status) {
  return { pending: 'orange', accepted: 'green', declined: 'red' }[status];
}