document.addEventListener("DOMContentLoaded", async function () {
  const tbody = document.getElementById("reviewsTableBody");
  const pendingEl = document.getElementById("ideasPendingCount");
  const approvedEl = document.getElementById("ideasApprovedCount");
  const rejectedEl = document.getElementById("ideasRejectedCount");
  const emptyState = document.getElementById("emptyState");

  const token = localStorage.getItem("authToken");
  if (!token) {
    window.location.href = "/login.html";
    return;
  }

  // State
  let allIdeas = [];
  let currentFilter = 'all';

  async function loadIdeas() {
    try {
      const res = await fetch("/api/faculty/reviews", {
        headers: { Authorization: "Bearer " + token }
      });

      if (!res.ok) {
        console.error("Failed to fetch reviews:", res.status);
        allIdeas = [];
        renderStats();
        renderTable();
        return;
      }

      const data = await res.json();
      allIdeas = Array.isArray(data) ? data : [];

      renderStats();
      renderTable();
    } catch (e) {
      console.error("Error loading ideas:", e);
      allIdeas = [];
      renderStats();
      renderTable();
    }
  }

  function renderStats() {
    let pending = 0, approved = 0, rejected = 0;

    allIdeas.forEach(i => {
      if (i.status === "admin_approved") pending++;
      if (i.status === "approved") approved++;
      if (i.status === "rejected") rejected++;
    });

    pendingEl.textContent = `${pending} pending`;
    approvedEl.textContent = `${approved} approved`;
    rejectedEl.textContent = `${rejected} rejected`;
  }

  function getFilteredIdeas() {
    if (currentFilter === 'all') {
      return allIdeas;
    } else if (currentFilter === 'pending') {
      return allIdeas.filter(i => i.status === 'admin_approved');
    } else if (currentFilter === 'approved') {
      return allIdeas.filter(i => i.status === 'approved');
    } else if (currentFilter === 'rejected') {
      return allIdeas.filter(i => i.status === 'rejected');
    }
    return allIdeas;
  }

  function renderTable() {
    tbody.innerHTML = "";
    const filteredIdeas = getFilteredIdeas();

    // Show/hide empty state
    if (emptyState) {
      emptyState.classList.toggle('hidden', filteredIdeas.length > 0);
    }

    filteredIdeas.forEach(idea => {
      const statusConfig = {
        admin_approved: {
          label: "Pending Faculty Review",
          pillClass: "bg-orange-50 text-orange-700 border-orange-200",
          icon: "fa-hourglass-half"
        },
        approved: {
          label: "Approved",
          pillClass: "bg-emerald-50 text-emerald-700 border-emerald-200",
          icon: "fa-check-circle"
        },
        rejected: {
          label: "Rejected",
          pillClass: "bg-rose-50 text-rose-700 border-rose-200",
          icon: "fa-times-circle"
        }
      };

      const config = statusConfig[idea.status] || {
        label: idea.status,
        pillClass: "bg-gray-50 text-gray-700 border-gray-200",
        icon: "fa-question-circle"
      };

      const isPending = idea.status === "admin_approved";

      const tr = document.createElement("tr");
      tr.className = "border-b border-gray-100 last:border-0 hover:bg-gray-50 transition-colors";

      tr.innerHTML = `
        <td class="py-3 px-4">
          <div class="font-semibold text-gray-900">${idea.title}</div>
          <div class="text-xs text-gray-500">
            ${new Date(idea.submitted_on).toLocaleDateString()}
          </div>
        </td>
        <td class="py-3 px-4 text-gray-700">${idea.student}</td>
        <td class="py-3 px-4">
          <span class="badge-pill ${config.pillClass} border">
            <i class="fas ${config.icon} text-xs"></i> ${config.label}
          </span>
        </td>
        <td class="py-3 px-4 text-right">
          <a href="faculty-idea.html?id=${idea.id}" 
             class="view-btn px-3 py-1 rounded-full border border-blue-300 text-blue-700 text-xs mr-2 hover:bg-blue-50 transition-colors inline-flex items-center gap-1">
            <i class="fas fa-eye"></i> View
          </a>
          ${isPending ? `
            <button class="approve-btn px-3 py-1 rounded-full border border-emerald-300 text-emerald-700 text-xs hover:bg-emerald-50 transition-colors" data-id="${idea.id}">
              <i class="fas fa-check"></i> Approve
            </button>
            <button class="reject-btn px-3 py-1 rounded-full border border-rose-300 text-rose-700 text-xs ml-2 hover:bg-rose-50 transition-colors" data-id="${idea.id}">
              <i class="fas fa-xmark"></i> Reject
            </button>
          ` : `
            <span class="text-xs text-gray-400">—</span>
          `}
        </td>
      `;

      tbody.appendChild(tr);
    });
  }

  // Handle approve/reject clicks
  tbody.addEventListener("click", async function (e) {
    const approveBtn = e.target.closest(".approve-btn");
    const rejectBtn = e.target.closest(".reject-btn");
    if (!approveBtn && !rejectBtn) return;

    const id = (approveBtn || rejectBtn).dataset.id;
    const decision = approveBtn ? "approved" : "rejected";

    const res = await fetch(`/api/faculty/reviews/${id}/decision`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        Authorization: "Bearer " + token
      },
      body: JSON.stringify({ decision })
    });

    if (!res.ok) {
      alert("Failed to update status");
      return;
    }

    // Update local state
    const idea = allIdeas.find(i => i.id === id);
    if (idea) idea.status = decision;

    renderStats();
    renderTable();
  });

  // Handle filter tab clicks
  function setupFilterTabs() {
    const filterBtns = document.querySelectorAll('[data-idea-filter]');

    filterBtns.forEach(btn => {
      btn.addEventListener('click', () => {
        // Update active state
        filterBtns.forEach(b => {
          b.classList.remove('bg-gray-900', 'text-white');
          b.classList.add('bg-gray-100', 'text-gray-700');
        });
        btn.classList.remove('bg-gray-100', 'text-gray-700');
        btn.classList.add('bg-gray-900', 'text-white');

        // Apply filter
        currentFilter = btn.dataset.ideaFilter;
        renderTable();
      });
    });
  }

  // Mobile nav toggle
  document.getElementById('mobile-menu-btn')?.addEventListener('click', function () {
    document.getElementById('mobile-menu')?.classList.toggle('hidden');
  });

  // Initialize
  setupFilterTabs();
  loadIdeas();
});