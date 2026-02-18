const token = localStorage.getItem('authToken');
if (!token) location.href = '/login.html';

// Fixed API endpoints to match backend routes
const API = {
  contents: '/api/contents',
  createContent: '/api/create-content',
  ideas: '/api/admin/submissions',  // Admin reviews student submissions
  faculty: '/api/admin/faculty'
};

const headers = {
  Authorization: `Bearer ${token}`,
  'Content-Type': 'application/json'
};

// Cache for content data (for editing)
let contentCache = [];

// Content type mappings
const ABOUT_TYPES = ['about_card', 'about_feature', 'about_stat', 'about_testimonial', 'team_member'];

// Map tabs to content types
const TAB_TO_TYPES = {
  resources: ['resource'],
  about: ABOUT_TYPES,
  events: ['event']
};

// Map content types to tabs
const TYPE_TO_TAB = {
  resource: 'resources',
  about_card: 'about',
  about_feature: 'about',
  about_stat: 'about',
  about_testimonial: 'about',
  team_member: 'about',
  event: 'events'
};

// Human readable names for content types
const TYPE_LABELS = {
  resource: 'Resource',
  about_card: 'About Card',
  about_feature: 'Feature',
  about_stat: 'Stat',
  about_testimonial: 'Testimonial',
  team_member: 'Team Member',
  event: 'Event'
};

// Field configuration for each content type
const TYPE_FIELD_CONFIG = {
  resource: {
    titleLabel: 'Title *',
    titlePlaceholder: 'Resource title',
    descLabel: 'Description',
    descPlaceholder: 'Brief description of the resource',
    imageLabel: 'File/Image URL',
    showImage: true,
    showIcon: false,
    showRole: false,
    showStatValue: false
  },
  about_card: {
    titleLabel: 'Card Title *',
    titlePlaceholder: 'e.g. Our Mission',
    descLabel: 'Card Description',
    descPlaceholder: 'Describe this section',
    imageLabel: 'Card Image URL',
    showImage: true,
    showIcon: false,
    showRole: false,
    showStatValue: false
  },
  about_feature: {
    titleLabel: 'Feature Title *',
    titlePlaceholder: 'e.g. Innovation Hub',
    descLabel: 'Feature Description',
    descPlaceholder: 'Describe this feature',
    imageLabel: 'Image URL (optional)',
    showImage: true,
    showIcon: true,
    showRole: false,
    showStatValue: false
  },
  about_stat: {
    titleLabel: 'Stat Label *',
    titlePlaceholder: 'e.g. Startups Incubated',
    descLabel: 'Additional Info (optional)',
    descPlaceholder: 'Any extra context',
    showImage: false,
    showIcon: false,
    showRole: false,
    showStatValue: true
  },
  about_testimonial: {
    titleLabel: 'Person Name *',
    titlePlaceholder: 'e.g. John Doe',
    descLabel: 'Testimonial Quote *',
    descPlaceholder: 'What they said about MIC...',
    imageLabel: 'Avatar/Photo URL',
    showImage: true,
    showIcon: false,
    showRole: true,
    showStatValue: false
  },
  team_member: {
    titleLabel: 'Member Name *',
    titlePlaceholder: 'e.g. Dr. Jane Smith',
    descLabel: 'Bio (optional)',
    descPlaceholder: 'Brief biography or expertise',
    imageLabel: 'Photo URL',
    showImage: true,
    showIcon: false,
    showRole: true,
    showStatValue: false
  },
  event: {
    titleLabel: 'Event Title *',
    titlePlaceholder: 'e.g. Startup Pitch Day',
    descLabel: 'Event Description',
    descPlaceholder: 'Describe the event',
    imageLabel: 'Event Banner URL',
    showImage: true,
    showIcon: false,
    showRole: false,
    showStatValue: false
  }
};

document.querySelectorAll('.tab-btn').forEach(btn => {
  btn.onclick = () => switchTab(btn.dataset.tab);
});

function switchTab(tab) {
  document.querySelectorAll('.section').forEach(s => s.classList.add('hidden'));
  document.getElementById(`${tab}-section`)?.classList.remove('hidden');

  document.querySelectorAll('.tab-btn').forEach(b => b.classList.remove('border-orange-500', 'border-b-2'));
  document.querySelector(`[data-tab="${tab}"]`)?.classList.add('border-orange-500', 'border-b-2');

  if (tab === 'ideas') loadIdeas();
  else if (tab === 'faculty') loadFaculty();
  else if (tab === 'work') {
    loadWork();
    loadCompanies();
  }
  else loadContent(tab);
}

async function loadContent(tab) {
  try {
    const res = await fetch(API.contents, { headers });
    if (!res.ok) throw new Error('Failed to fetch content');
    const data = await res.json();
    contentCache = data || [];

    // Get all types for this tab
    const typesForTab = TAB_TO_TYPES[tab] || [];
    const list = contentCache.filter(i => typesForTab.includes(i.content_type));
    const listEl = document.getElementById(`${tab}-list`);

    if (list.length === 0) {
      listEl.innerHTML = '<p class="text-gray-500">No items found. Add one using the button above.</p>';
      return;
    }

    listEl.innerHTML = list.map(c => {
      // Parse content_data if it's a string
      let contentData = c.content_data;
      if (typeof contentData === 'string') {
        try { contentData = JSON.parse(contentData); } catch (e) { contentData = {}; }
      }

      return `
        <div class="bg-white p-4 rounded shadow">
          <div class="flex justify-between items-start mb-2">
            <span class="text-xs bg-gray-200 text-gray-700 px-2 py-1 rounded">${escapeHtml(TYPE_LABELS[c.content_type] || c.content_type)}</span>
            <span class="text-xs text-gray-400">#${c.order_index || 0}</span>
          </div>
          ${c.content_type === 'about_stat' && contentData?.stat_value ? `<p class="text-2xl font-bold text-orange-500">${escapeHtml(contentData.stat_value)}</p>` : ''}
          <h3 class="font-bold">${escapeHtml(c.title)}</h3>
          ${contentData?.role ? `<p class="text-sm text-orange-600">${escapeHtml(contentData.role)}</p>` : ''}
          <p class="text-sm text-gray-600 mt-1">${escapeHtml(c.description || '')}</p>
          ${contentData?.icon ? `<p class="text-lg mt-1">${escapeHtml(contentData.icon)}</p>` : ''}
          ${c.image_url ? `<img src="${escapeHtml(c.image_url)}" class="w-16 h-16 object-cover mt-2 rounded"/>` : ''}
          <div class="flex gap-3 mt-3 pt-2 border-t">
            <button onclick="editContent('${c.id}')" class="text-blue-600 hover:underline text-sm">Edit</button>
            <button onclick="deleteContent('${c.id}')" class="text-red-600 hover:underline text-sm">Delete</button>
          </div>
        </div>
      `;
    }).join('');
  } catch (err) {
    console.error('Error loading content:', err);
    document.getElementById(`${tab}-list`).innerHTML = '<p class="text-red-500">Error loading content. Please try again.</p>';
  }
}

// Simple HTML escape to prevent XSS
function escapeHtml(str) {
  if (!str) return '';
  return String(str).replace(/&/g, '&amp;')
    .replace(/</g, '&lt;')
    .replace(/>/g, '&gt;')
    .replace(/"/g, '&quot;')
    .replace(/'/g, '&#039;');
}

// Update add button label when dropdown changes
window.updateAddButtonLabel = function () {
  const type = document.getElementById('addAboutType').value;
  document.getElementById('addAboutBtn').textContent = `+ Add ${TYPE_LABELS[type] || type}`;
};

// Configure form fields based on content type
function configureFormFields(type) {
  const config = TYPE_FIELD_CONFIG[type] || TYPE_FIELD_CONFIG.resource;

  // Update labels
  document.getElementById('titleLabel').textContent = config.titleLabel;
  document.getElementById('title').placeholder = config.titlePlaceholder;
  document.getElementById('descriptionLabel').textContent = config.descLabel;
  document.getElementById('description').placeholder = config.descPlaceholder;

  // Show/hide fields
  document.getElementById('fieldImageUrl').classList.toggle('hidden', !config.showImage);
  if (config.imageLabel) {
    document.getElementById('imageLabel').textContent = config.imageLabel;
  }

  document.getElementById('fieldIcon').classList.toggle('hidden', !config.showIcon);
  document.getElementById('fieldRole').classList.toggle('hidden', !config.showRole);
  document.getElementById('fieldStatValue').classList.toggle('hidden', !config.showStatValue);
  document.getElementById('eventFields').classList.toggle('hidden', type !== 'event');
}

window.openModal = function (type, id = null) {
  modal.classList.remove('hidden');
  contentType.value = type;
  contentId.value = id || '';
  modalTitle.innerText = id ? `Edit ${TYPE_LABELS[type] || type}` : `Add ${TYPE_LABELS[type] || type}`;

  // Configure fields for this type
  configureFormFields(type);

  // Reset form first
  contentForm.reset();

  // Pre-populate form if editing
  if (id) {
    const item = contentCache.find(c => c.id === id);
    if (item) {
      title.value = item.title || '';
      description.value = item.description || '';
      imageUrl.value = item.image_url || '';
      orderIndex.value = item.order_index || 0;

      // Parse content_data
      let contentData = item.content_data;
      if (typeof contentData === 'string') {
        try { contentData = JSON.parse(contentData); } catch (e) { contentData = {}; }
      }
      contentData = contentData || {};

      // Populate type-specific fields
      if (document.getElementById('iconField')) {
        document.getElementById('iconField').value = contentData.icon || '';
      }
      if (document.getElementById('roleField')) {
        document.getElementById('roleField').value = contentData.role || '';
      }
      if (document.getElementById('statValue')) {
        document.getElementById('statValue').value = contentData.stat_value || '';
      }

      if (type === 'event') {
        eventDate.value = contentData.event_date || '';
        venue.value = contentData.venue || '';
        price.value = contentData.price || '';
        regLink.value = contentData.registration_link || '';
      }
    }
  }
};

window.closeModal = function () {
  modal.classList.add('hidden');
  contentForm.reset();
  contentId.value = '';
};

window.editContent = function (id) {
  const item = contentCache.find(c => c.id === id);
  if (item) {
    openModal(item.content_type, id);
  }
};

window.deleteContent = async function (id) {
  if (!confirm('Are you sure you want to delete this item?')) return;

  try {
    const res = await fetch(`${API.contents}/${id}`, {
      method: 'DELETE',
      headers
    });

    if (!res.ok) throw new Error('Failed to delete');

    // Reload current tab
    const activeTab = document.querySelector('.tab-btn.border-orange-500')?.dataset.tab || 'resources';
    loadContent(activeTab);
  } catch (err) {
    console.error('Error deleting content:', err);
    alert('Failed to delete content. Please try again.');
  }
};

contentForm.onsubmit = async e => {
  e.preventDefault();

  const type = contentType.value;

  // Build content_data based on type
  const contentData = {};

  if (type === 'about_feature') {
    const icon = document.getElementById('iconField')?.value;
    if (icon) contentData.icon = icon;
  }

  if (type === 'about_testimonial' || type === 'team_member') {
    const role = document.getElementById('roleField')?.value;
    if (role) contentData.role = role;
  }

  if (type === 'about_stat') {
    const statVal = document.getElementById('statValue')?.value;
    if (statVal) contentData.stat_value = statVal;
  }

  if (type === 'event') {
    contentData.event_date = eventDate.value;
    contentData.venue = venue.value;
    contentData.price = price.value;
    contentData.registration_link = regLink.value;
  }

  const payload = {
    content_type: type,
    title: title.value,
    description: description.value || null,
    image_url: imageUrl.value || null,
    order_index: Number(orderIndex.value || 0),
    is_active: true,
    content_data: contentData
  };

  try {
    const isEdit = !!contentId.value;
    const url = isEdit ? `${API.contents}/${contentId.value}` : API.createContent;
    const method = isEdit ? 'PUT' : 'POST';

    const res = await fetch(url, {
      method,
      headers,
      body: JSON.stringify(payload)
    });

    if (!res.ok) {
      const errText = await res.text();
      throw new Error(errText || 'Failed to save content');
    }

    closeModal();
    // Reload correct tab based on content type
    const targetTab = TYPE_TO_TAB[type] || 'resources';
    switchTab(targetTab);
  } catch (err) {
    console.error('Error saving content:', err);
    alert('Failed to save content: ' + err.message);
  }
};

// Ideas/Submissions Management - Admin reviews and manages ideas
let ideasCache = [];
let allFacultyCache = [];
let ideasFilter = 'all';

async function loadIdeas() {
  try {
    // Fetch ALL ideas (not just pending)
    const res = await fetch('/api/admin/submissions/all', { headers });
    if (!res.ok) throw new Error('Failed to fetch ideas');
    ideasCache = await res.json() || [];

    // Also load faculty for assignment
    await loadAllFacultyForAssignment();

    renderIdeas();
  } catch (err) {
    console.error('Error loading ideas:', err);
    document.getElementById('ideas-list').innerHTML = '<p class="text-red-500">Error loading submissions. You may not have admin permissions.</p>';
  }
}

async function loadAllFacultyForAssignment() {
  try {
    const res = await fetch(API.faculty, { headers });
    if (res.ok) {
      allFacultyCache = await res.json() || [];
    }
  } catch (err) {
    console.error('Error loading faculty:', err);
    allFacultyCache = [];
  }
}

function renderIdeas() {
  const ideasList = document.getElementById('ideas-list');

  // Filter ideas based on current filter
  let filtered = ideasCache;
  if (ideasFilter !== 'all') {
    filtered = ideasCache.filter(i => i.status === ideasFilter);
  }

  if (!filtered || filtered.length === 0) {
    ideasList.innerHTML = '<p class="text-gray-500">No submissions found.</p>';
    return;
  }

  const statusColors = {
    submitted: 'bg-yellow-100 text-yellow-700',
    admin_approved: 'bg-green-100 text-green-700',
    admin_rejected: 'bg-red-100 text-red-700',
    approved: 'bg-blue-100 text-blue-700',
    rejected: 'bg-red-100 text-red-700'
  };

  const statusLabels = {
    submitted: 'Pending Review',
    admin_approved: 'Approved for Faculty',
    admin_rejected: 'Rejected',
    approved: 'Faculty Approved',
    rejected: 'Faculty Rejected'
  };

  ideasList.innerHTML = filtered.map(i => {
    const statusColor = statusColors[i.status] || 'bg-gray-100 text-gray-700';
    const statusLabel = statusLabels[i.status] || i.status;
    const tags = i.tags && i.tags.length > 0 ? i.tags : [];
    const isPending = i.status === 'submitted';

    return `
      <div class="bg-white p-4 rounded shadow">
        <div class="flex justify-between items-start mb-2">
          <h3 class="font-bold text-lg">${escapeHtml(i.title)}</h3>
          <span class="text-xs px-2 py-1 rounded ${statusColor}">${statusLabel}</span>
        </div>
        <p class="text-sm text-gray-700 mt-1">${escapeHtml(i.description || 'No description')}</p>
        <p class="text-sm text-gray-600 mt-2">Submitted by: <span class="font-medium">${escapeHtml(i.student || 'Unknown')}</span></p>
        <p class="text-xs text-gray-400">Submitted: ${new Date(i.submitted_on).toLocaleDateString()}</p>
        ${i.file_path ? `<p class="text-xs text-blue-600 mt-1"><a href="${escapeHtml(i.file_path)}" target="_blank">📄 View Attached File</a></p>` : ''}
        
        <!-- Domain & Tags -->
        <div class="mt-2 flex flex-wrap gap-1">
          ${i.domain ? `<span class="text-xs px-2 py-0.5 rounded bg-purple-100 text-purple-700">${escapeHtml(i.domain)}</span>` : ''}
          ${tags.map(t => `<span class="text-xs px-2 py-0.5 rounded bg-gray-100 text-gray-600">${escapeHtml(t)}</span>`).join('')}
        </div>
        
        <!-- Action Buttons -->
        <div class="flex flex-wrap gap-2 mt-4 pt-3 border-t">
          ${isPending ? `
            <button onclick="decide('${i.id}','approved')" class="bg-green-500 text-white px-3 py-1.5 rounded text-sm hover:bg-green-600">✓ Approve</button>
            <button onclick="decide('${i.id}','rejected')" class="bg-red-500 text-white px-3 py-1.5 rounded text-sm hover:bg-red-600">✗ Reject</button>
          ` : ''}
          <button onclick="openFacultyAssignModal('${i.id}')" class="bg-blue-500 text-white px-3 py-1.5 rounded text-sm hover:bg-blue-600">👥 Assign Faculty</button>
          <button onclick="openTagsModal('${i.id}')" class="bg-purple-500 text-white px-3 py-1.5 rounded text-sm hover:bg-purple-600">🏷️ Tags</button>
        </div>
      </div>
    `;
  }).join('');
}

// Setup ideas filter tabs
function setupIdeasFilterTabs() {
  const filterBtns = document.querySelectorAll('[data-ideas-filter]');
  filterBtns.forEach(btn => {
    btn.addEventListener('click', () => {
      // Update button styles
      filterBtns.forEach(b => {
        b.classList.remove('bg-gray-900', 'text-white');
        b.classList.add('bg-gray-100', 'text-gray-700');
      });
      btn.classList.remove('bg-gray-100', 'text-gray-700');
      btn.classList.add('bg-gray-900', 'text-white');

      // Apply filter
      ideasFilter = btn.dataset.ideasFilter;
      renderIdeas();
    });
  });
}

// Initialize filter tabs after DOM load
document.addEventListener('DOMContentLoaded', setupIdeasFilterTabs);

window.decide = async (id, decision) => {
  try {
    const res = await fetch(`${API.ideas}/${id}/decision`, {
      method: 'POST',
      headers,
      body: JSON.stringify({ decision })
    });

    if (!res.ok) throw new Error('Failed to update decision');
    loadIdeas();
  } catch (err) {
    console.error('Error making decision:', err);
    alert('Failed to update decision. Please try again.');
  }
};

// =====================
// FACULTY ASSIGNMENT
// =====================

window.openFacultyAssignModal = async function (ideaId) {
  const idea = ideasCache.find(i => i.id === ideaId);
  if (!idea) return;

  document.getElementById('assignIdeaId').value = ideaId;
  document.getElementById('assignIdeaTitle').textContent = `For: ${idea.title}`;

  // Populate faculty dropdown
  const select = document.getElementById('facultyToAssign');
  select.innerHTML = '<option value="">Select a faculty member...</option>' +
    allFacultyCache.map(f => `<option value="${f.id}">${escapeHtml(f.name)} (${escapeHtml(f.email)})</option>`).join('');

  // Load currently assigned faculty
  await loadAssignedFaculty(ideaId);

  document.getElementById('facultyAssignModal').classList.remove('hidden');
};

async function loadAssignedFaculty(ideaId) {
  const list = document.getElementById('assignedFacultyList');
  list.innerHTML = '<p class="text-xs text-gray-400">Loading...</p>';

  try {
    const res = await fetch(`/api/admin/submissions/${ideaId}/faculty`, { headers });
    if (!res.ok) throw new Error('Failed to load');

    const assigned = await res.json() || [];

    if (assigned.length === 0) {
      list.innerHTML = '<p class="text-xs text-gray-400">No faculty assigned yet</p>';
      return;
    }

    list.innerHTML = assigned.map(f => `
      <div class="flex items-center justify-between bg-gray-50 px-3 py-2 rounded">
        <span class="text-sm">${escapeHtml(f.faculty_name)}</span>
        <button onclick="removeFacultyAssignment('${ideaId}', '${f.faculty_id}')" class="text-red-500 hover:text-red-700 text-sm">Remove</button>
      </div>
    `).join('');
  } catch (err) {
    console.error('Error loading assigned faculty:', err);
    list.innerHTML = '<p class="text-xs text-red-500">Error loading assigned faculty</p>';
  }
}

window.addFacultyAssignment = async function () {
  const ideaId = document.getElementById('assignIdeaId').value;
  const facultyId = document.getElementById('facultyToAssign').value;

  if (!facultyId) {
    alert('Please select a faculty member');
    return;
  }

  try {
    const res = await fetch(`/api/admin/submissions/${ideaId}/assign-faculty`, {
      method: 'POST',
      headers,
      body: JSON.stringify({ faculty_id: facultyId })
    });

    if (!res.ok) throw new Error('Failed to assign');

    // Reload assigned list
    await loadAssignedFaculty(ideaId);
    document.getElementById('facultyToAssign').value = '';
  } catch (err) {
    console.error('Error assigning faculty:', err);
    alert('Failed to assign faculty. Please try again.');
  }
};

window.removeFacultyAssignment = async function (ideaId, facultyId) {
  if (!confirm('Remove this faculty assignment?')) return;

  try {
    const res = await fetch(`/api/admin/submissions/${ideaId}/assign-faculty/${facultyId}`, {
      method: 'DELETE',
      headers
    });

    if (!res.ok) throw new Error('Failed to remove');

    await loadAssignedFaculty(ideaId);
  } catch (err) {
    console.error('Error removing faculty:', err);
    alert('Failed to remove faculty. Please try again.');
  }
};

window.closeFacultyAssignModal = function () {
  document.getElementById('facultyAssignModal').classList.add('hidden');
};

// =====================
// TAGS MANAGEMENT
// =====================

window.openTagsModal = function (ideaId) {
  const idea = ideasCache.find(i => i.id === ideaId);
  if (!idea) return;

  document.getElementById('tagsIdeaId').value = ideaId;
  document.getElementById('tagsIdeaTitle').textContent = `For: ${idea.title}`;
  document.getElementById('ideaDomain').value = idea.domain || '';
  document.getElementById('ideaTags').value = (idea.tags || []).join(', ');

  document.getElementById('tagsModal').classList.remove('hidden');
};

window.closeTagsModal = function () {
  document.getElementById('tagsModal').classList.add('hidden');
};

document.getElementById('tagsForm')?.addEventListener('submit', async (e) => {
  e.preventDefault();

  const ideaId = document.getElementById('tagsIdeaId').value;
  const domain = document.getElementById('ideaDomain').value.trim();
  const tagsInput = document.getElementById('ideaTags').value;
  const tags = tagsInput.split(',').map(t => t.trim()).filter(t => t.length > 0);

  try {
    const res = await fetch(`/api/admin/submissions/${ideaId}/tags`, {
      method: 'PUT',
      headers,
      body: JSON.stringify({ domain, tags })
    });

    if (!res.ok) throw new Error('Failed to save tags');

    closeTagsModal();
    loadIdeas(); // Reload to show updated tags
  } catch (err) {
    console.error('Error saving tags:', err);
    alert('Failed to save tags. Please try again.');
  }
});

// Faculty Management
async function loadFaculty() {
  try {
    const res = await fetch(API.faculty, { headers });
    if (!res.ok) throw new Error('Failed to fetch faculty');
    const data = await res.json();

    const facultyList = document.getElementById('faculty-list');

    if (!data || data.length === 0) {
      facultyList.innerHTML = '<p class="text-gray-500">No faculty members found.</p>';
      return;
    }

    facultyList.innerHTML = data.map(f => `
      <div class="bg-white p-4 rounded shadow flex justify-between items-center">
        <div>
          <h3 class="font-bold">${escapeHtml(f.name)}</h3>
          <p class="text-sm text-gray-600">${escapeHtml(f.email)}</p>
        </div>
        <div class="flex gap-2">
          <button onclick="openFacultyModal('${f.id}', '${escapeHtml(f.name)}', '${escapeHtml(f.email)}')" class="text-blue-600 hover:underline">Edit</button>
          <button onclick="removeFaculty('${f.id}')" class="text-red-600 hover:underline">Remove</button>
        </div>
      </div>
    `).join('');
  } catch (err) {
    console.error('Error loading faculty:', err);
    document.getElementById('faculty-list').innerHTML = '<p class="text-red-500">Error loading faculty. You may not have admin permissions.</p>';
  }
}

window.removeFaculty = async function (id) {
  if (!confirm('Are you sure you want to remove this faculty member?')) return;

  try {
    const res = await fetch(`${API.faculty}/${id}`, {
      method: 'DELETE',
      headers
    });

    if (!res.ok) throw new Error('Failed to remove faculty');
    loadFaculty();
  } catch (err) {
    console.error('Error removing faculty:', err);
    alert('Failed to remove faculty. Please try again.');
  }
};

// Faculty Modal Handlers
window.openAddFacultyModal = function () {
  document.getElementById('facultyModal').classList.remove('hidden');
  document.getElementById('facultyModalTitle').innerText = 'Add Faculty';
  document.getElementById('facultyForm').reset();
  document.getElementById('facultyId').value = '';
  document.getElementById('facultyPasswordField').classList.remove('hidden');
  document.getElementById('facultyPassword').required = true;
};

window.openFacultyModal = function (id, name, email) {
  document.getElementById('facultyModal').classList.remove('hidden');
  document.getElementById('facultyModalTitle').innerText = 'Edit Faculty';
  document.getElementById('facultyId').value = id;
  document.getElementById('facultyName').value = name;
  document.getElementById('facultyEmail').value = email;
  document.getElementById('facultyPasswordField').classList.add('hidden');
  document.getElementById('facultyPassword').required = false;
};

window.closeFacultyModal = function () {
  document.getElementById('facultyModal').classList.add('hidden');
  document.getElementById('facultyForm').reset();
};

document.getElementById('facultyForm').onsubmit = async (e) => {
  e.preventDefault();

  const id = document.getElementById('facultyId').value;
  const name = document.getElementById('facultyName').value;
  const email = document.getElementById('facultyEmail').value;
  const password = document.getElementById('facultyPassword').value;

  const payload = { name, email };
  if (password) {
    payload.password = password;
  }

  try {
    const isEdit = !!id;
    const url = isEdit ? `${API.faculty}/${id}` : API.faculty;
    const method = isEdit ? 'PUT' : 'POST';

    const res = await fetch(url, {
      method,
      headers,
      body: JSON.stringify(payload)
    });

    if (!res.ok) {
      const error = await res.text();
      throw new Error(error || 'Failed to save faculty');
    }

    closeFacultyModal();
    loadFaculty();
  } catch (err) {
    console.error('Error saving faculty:', err);
    alert('Failed to save faculty: ' + err.message);
  }
};

// =====================
// WORK PIPELINE MANAGEMENT
// =====================

let companiesCache = [];
let workCache = [];

async function loadWork() {
  try {
    const res = await fetch('/api/admin/work', { headers });
    if (!res.ok) throw new Error('Failed to fetch work items');
    workCache = await res.json() || [];

    const workList = document.getElementById('work-list');
    if (!workCache || workCache.length === 0) {
      workList.innerHTML = '<p class="text-gray-500">No incubation projects found.</p>';
      return;
    }

    const stageLabels = {
      under_incubation: 'Under Incubation',
      looking_for_funding: 'Looking for Funding',
      funded: 'Found Companies'
    };

    const stageColors = {
      under_incubation: 'bg-blue-100 text-blue-700',
      looking_for_funding: 'bg-yellow-100 text-yellow-700',
      funded: 'bg-green-100 text-green-700'
    };

    workList.innerHTML = workCache.map(w => `
      <div class="bg-white p-4 rounded shadow">
        <div class="flex justify-between items-start mb-2">
          <h3 class="font-bold text-lg">${escapeHtml(w.title)}</h3>
          <span class="text-xs px-2 py-1 rounded ${stageColors[w.stage] || 'bg-gray-100'}">
            ${stageLabels[w.stage] || w.stage}
          </span>
        </div>
        <p class="text-sm text-gray-700 mt-1">${escapeHtml(w.description || 'No description')}</p>
        <div class="mt-2 text-xs text-gray-500">
          <p>Progress: ${w.progress_percent || 0}%</p>
          ${w.company_name ? `<p>Company: ${escapeHtml(w.company_name)}</p>` : ''}
        </div>
        <div class="flex gap-2 mt-3 pt-3 border-t">
          <button onclick="openWorkModal('${w.id}')" class="text-blue-600 text-sm hover:underline">Edit</button>
          <button onclick="deleteWork('${w.id}')" class="text-red-600 text-sm hover:underline">Delete</button>
        </div>
      </div>
    `).join('');
  } catch (err) {
    console.error('Error loading work:', err);
    document.getElementById('work-list').innerHTML = '<p class="text-red-500">Error loading work items.</p>';
  }
}

async function loadCompanies() {
  try {
    const res = await fetch('/api/admin/companies', { headers });
    if (!res.ok) throw new Error('Failed to fetch companies');
    companiesCache = await res.json() || [];

    const companiesList = document.getElementById('companies-list');
    if (!companiesCache || companiesCache.length === 0) {
      companiesList.innerHTML = '<p class="text-gray-500">No companies added yet.</p>';
      return;
    }

    companiesList.innerHTML = companiesCache.map(c => `
      <div class="bg-white p-3 rounded shadow flex items-center gap-3">
        ${c.logo_url ? `<img src="${escapeHtml(c.logo_url)}" class="w-10 h-10 object-contain rounded" onerror="this.style.display='none'" />` : ''}
        <div class="flex-1">
          <p class="font-medium">${escapeHtml(c.name)}</p>
        </div>
        <button onclick="deleteCompany('${c.id}')" class="text-red-500 hover:text-red-700">
          <span class="text-lg">×</span>
        </button>
      </div>
    `).join('');

    // Also update any company dropdowns
    updateCompanyDropdowns();
  } catch (err) {
    console.error('Error loading companies:', err);
    document.getElementById('companies-list').innerHTML = '<p class="text-red-500">Error loading companies.</p>';
  }
}

function updateCompanyDropdowns() {
  const workCompanySelect = document.getElementById('workCompany');
  if (workCompanySelect) {
    const currentValue = workCompanySelect.value;
    workCompanySelect.innerHTML = '<option value="">None</option>' +
      (companiesCache || []).map(c => `<option value="${c.id}">${escapeHtml(c.name)}</option>`).join('');
    workCompanySelect.value = currentValue;
  }
}

window.openAddCompanyModal = function () {
  document.getElementById('companyForm').reset();
  document.getElementById('companyModal').classList.remove('hidden');
};

window.closeCompanyModal = function () {
  document.getElementById('companyModal').classList.add('hidden');
};

document.getElementById('companyForm').onsubmit = async (e) => {
  e.preventDefault();

  const name = document.getElementById('companyName').value;
  const logoUrl = document.getElementById('companyLogo').value || null;

  try {
    const res = await fetch('/api/admin/companies', {
      method: 'POST',
      headers,
      body: JSON.stringify({ name, logo_url: logoUrl })
    });

    if (!res.ok) throw new Error('Failed to add company');

    closeCompanyModal();
    loadCompanies();
  } catch (err) {
    console.error('Error adding company:', err);
    alert('Failed to add company: ' + err.message);
  }
};

window.deleteCompany = async function (id) {
  if (!confirm('Delete this company?')) return;

  try {
    const res = await fetch(`/api/admin/companies/${id}`, {
      method: 'DELETE',
      headers
    });

    if (!res.ok) throw new Error('Failed to delete company');
    loadCompanies();
    loadWork(); // Refresh work list in case company was associated
  } catch (err) {
    console.error('Error deleting company:', err);
    alert('Failed to delete company: ' + err.message);
  }
};

window.openWorkModal = function (id) {
  const work = workCache.find(w => w.id === id);
  if (!work) return;

  document.getElementById('workId').value = id;
  document.getElementById('workTitle').value = work.title;
  document.getElementById('workDescription').value = work.description || '';
  document.getElementById('workStage').value = work.stage;
  document.getElementById('workProgress').value = work.progress_percent || 0;
  document.getElementById('workCompany').value = work.company_id || '';

  updateCompanyDropdowns();
  document.getElementById('workModal').classList.remove('hidden');
};

window.closeWorkModal = function () {
  document.getElementById('workModal').classList.add('hidden');
  document.getElementById('workForm').reset();
};

document.getElementById('workForm').onsubmit = async (e) => {
  e.preventDefault();

  const id = document.getElementById('workId').value;
  const title = document.getElementById('workTitle').value;
  const description = document.getElementById('workDescription').value;
  const stage = document.getElementById('workStage').value;
  const progressPercent = parseInt(document.getElementById('workProgress').value) || 0;
  const companyId = document.getElementById('workCompany').value || null;

  try {
    const res = await fetch(`/api/admin/work/${id}`, {
      method: 'PUT',
      headers,
      body: JSON.stringify({
        title,
        description,
        stage,
        progress_percent: progressPercent,
        company_id: companyId
      })
    });

    if (!res.ok) throw new Error('Failed to update work item');

    closeWorkModal();
    loadWork();
  } catch (err) {
    console.error('Error updating work:', err);
    alert('Failed to update work item: ' + err.message);
  }
};

window.deleteWork = async function (id) {
  if (!confirm('Delete this work item? This cannot be undone.')) return;

  try {
    const res = await fetch(`/api/admin/work/${id}`, {
      method: 'DELETE',
      headers
    });

    if (!res.ok) throw new Error('Failed to delete work item');
    loadWork();
  } catch (err) {
    console.error('Error deleting work:', err);
    alert('Failed to delete work item: ' + err.message);
  }
};

// Initialize
switchTab('resources');