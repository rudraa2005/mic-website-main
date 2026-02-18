document.addEventListener('DOMContentLoaded', async function () {
  const listEl = document.getElementById('eventsList');
  const acceptedEl = document.getElementById('eventsAcceptedCount');
  const pendingEl = document.getElementById('eventsPendingCount');
  const declinedEl = document.getElementById('eventsDeclinedCount');
  const filterButtons = document.querySelectorAll('[data-filter]');

  const token = localStorage.getItem('authToken');
  if (!token) {
    window.location.href = '/login.html';
    return;
  }

  let invitations = []; // single source of truth
  let activeFilter = 'all';

  // ----------------------------
  // Fetch data
  // ----------------------------
  async function loadInvitations() {
    const res = await fetch('/api/faculty/events/invitations', {
      headers: {
        Authorization: 'Bearer ' + token,
      },
    });

    if (!res.ok) {
      alert('Failed to load invitations');
      return;
    }

    invitations = await res.json();
    renderStats();
    renderEvents();
  }

  // ----------------------------
  // Stats
  // ----------------------------
  function renderStats() {
    let accepted = 0;
    let pending = 0;
    let declined = 0;

    invitations.forEach(inv => {
      if (inv.status === 'accepted') accepted++;
      if (inv.status === 'pending') pending++;
      if (inv.status === 'declined') declined++;
    });

    acceptedEl.textContent = `${accepted} accepted`;
    pendingEl.textContent = `${pending} pending`;
    declinedEl.textContent = `${declined} declined`;
  }

  // ----------------------------
  // Events list
  // ----------------------------
  function renderEvents() {
    const now = new Date();
    listEl.innerHTML = '';

    invitations.forEach(inv => {
      const eventDate = new Date(inv.event_date + 'T00:00:00');
      const isPast = eventDate < now;

      if (activeFilter === 'upcoming' && isPast) return;
      if (activeFilter === 'past' && !isPast) return;
      if (activeFilter === 'pending' && inv.status !== 'pending') return;

      const statusLabel = {
        pending: 'Awaiting your RSVP',
        accepted: 'You have accepted',
        declined: 'You have declined',
      }[inv.status];

      const statusPillClass = {
        pending: 'badge-pill bg-orange-100 text-orange-800 border border-orange-300',
        accepted: 'badge-pill bg-emerald-100 text-emerald-800 border border-emerald-300',
        declined: 'badge-pill bg-rose-100 text-rose-800 border border-rose-300',
      }[inv.status];

      const statusIcon = {
        pending: 'fa-bell',
        accepted: 'fa-circle-check',
        declined: 'fa-circle-xmark',
      }[inv.status];

      const card = document.createElement('div');
      card.className =
        'rounded-2xl border border-gray-200 bg-white px-4 py-4 md:px-5 md:py-5 flex flex-col md:flex-row md:items-start md:justify-between gap-4';
      card.dataset.invitationId = inv.invitation_id;

      card.innerHTML = `
        <div class="flex-1">
          <div class="flex items-center justify-between mb-1">
            <h3 class="font-semibold text-gray-900 md:text-base">${inv.title}</h3>
            <span class="${statusPillClass} text-[11px]">
              <i class="fas ${statusIcon} text-[10px]"></i>
              ${statusLabel}
            </span>
          </div>
          <p class="text-xs text-gray-600 mb-2 flex items-center gap-2">
            <span class="inline-flex items-center gap-1 px-2 py-0.5 rounded-full bg-gray-100 text-gray-700">
              <i class="fas fa-location-dot text-[10px]"></i>${inv.venue}
            </span>
            <span class="inline-flex items-center gap-1 px-2 py-0.5 rounded-full bg-gray-100 text-gray-700">
              ${inv.price}
            </span>
          </p>
          <div class="flex flex-wrap items-center gap-x-4 gap-y-1 text-[11px] text-gray-500">
            <span><i class="far fa-calendar"></i> ${inv.event_date}</span>
            <span><i class="far fa-clock"></i> ${isPast ? 'Past event' : 'Upcoming event'}</span>
          </div>
        </div>

        <div class="flex md:flex-col items-center md:items-end gap-2 text-xs">
          <button class="rsvp-accept px-3 py-1.5 rounded-full bg-emerald-600 text-white hover:bg-emerald-700 disabled:opacity-40"
            ${inv.status !== 'pending' ? 'disabled' : ''}>
            <i class="fas fa-check"></i> Accept
          </button>
          <button class="rsvp-decline px-3 py-1.5 rounded-full bg-rose-600 text-white hover:bg-rose-700 disabled:opacity-40"
            ${inv.status !== 'pending' ? 'disabled' : ''}>
            <i class="fas fa-xmark"></i> Decline
          </button>
        </div>
      `;

      listEl.appendChild(card);
    });
  }

  // ----------------------------
  // RSVP actions
  // ----------------------------
  listEl.addEventListener('click', async function (e) {
    const acceptBtn = e.target.closest('.rsvp-accept');
    const declineBtn = e.target.closest('.rsvp-decline');
    if (!acceptBtn && !declineBtn) return;

    const card = e.target.closest('[data-invitation-id]');
    if (!card) return;

    const invitationId = card.dataset.invitationId;
    const status = acceptBtn ? 'accepted' : 'declined';

    const res = await fetch(
      `/api/faculty/events/invitations/${invitationId}/rsvp`,
      {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          Authorization: 'Bearer ' + token,
        },
        body: JSON.stringify({ status }),
      }
    );

    if (!res.ok) {
      alert('Failed to update RSVP');
      return;
    }

    // update local state
    const inv = invitations.find(i => i.invitation_id === invitationId);
    if (inv) inv.status = status;

    renderStats();
    renderEvents();
  });

  // ----------------------------
  // Filters
  // ----------------------------
  filterButtons.forEach(btn => {
    btn.addEventListener('click', function () {
      filterButtons.forEach(b => {
        b.classList.remove('bg-gray-900', 'text-white');
        b.classList.add('bg-gray-100', 'text-gray-700');
      });
      this.classList.add('bg-gray-900', 'text-white');
      this.classList.remove('bg-gray-100', 'text-gray-700');

      activeFilter = this.getAttribute('data-filter') || 'all';
      renderEvents();
    });
  });

  document.getElementById('mobile-menu-btn')?.addEventListener('click', () => {
    document.getElementById('mobile-menu')?.classList.toggle('hidden');
  });

  // init
  await loadInvitations();
});