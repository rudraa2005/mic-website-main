document.addEventListener('DOMContentLoaded', async function () {
  const backBtn = document.getElementById('backButton');
  const ideaNotFound = document.getElementById('ideaNotFound');
  const ideaContainer = document.getElementById('ideaContainer');

  const titleEl = document.getElementById('ideaTitle');
  const studentEl = document.getElementById('ideaStudent');
  const submittedOnEl = document.getElementById('ideaSubmittedOn');
  const statusBadge = document.getElementById('ideaStatusBadge');
  const domainBadge = document.getElementById('ideaDomainBadge');
  const descriptionEl = document.getElementById('ideaDescription');
  const studentEmailEl = document.getElementById('ideaStudentEmail');
  const attachmentEl = document.getElementById('ideaAttachment');
  const viewSubmissionBtn = document.getElementById('viewSubmissionBtn');

  const token = localStorage.getItem('authToken');
  if (!token) {
    window.location.href = '/login.html';
    return;
  }

  function getSubmissionId() {
    return new URLSearchParams(window.location.search).get('id');
  }

  function getMode() {
    return new URLSearchParams(window.location.search).get('mode');
  }

  async function fetchIdea(submissionId) {
    try {
      const res = await fetch(`/api/faculty/reviews/${submissionId}`, {
        headers: {
          Authorization: 'Bearer ' + token
        }
      });

      if (!res.ok) {
        const errorText = await res.text().catch(() => 'No error details');
        console.error(`Fetch error: ${res.status} ${res.statusText} - ${errorText}`);
        return { error: true, status: res.status, message: errorText };
      }
      return res.json();
    } catch (e) {
      console.error('Network or parsing error:', e);
      return { error: true, status: 0, message: e.message };
    }
  }

  function renderIdea(idea) {
    titleEl.textContent = idea.title;
    studentEl.textContent = `Submitted by ${idea.student}`;
    studentEmailEl.textContent = idea.email;
    descriptionEl.textContent = idea.description || 'No description provided.';

    submittedOnEl.textContent =
      `Submitted on ${new Date(idea.submitted_on).toLocaleDateString()}`;

    if (idea.file_path) {
      // Extract filename from path (e.g., ./uploads/uuid_filename.ext -> filename.ext)
      const parts = idea.file_path.split('_');
      const filename = parts.length > 1 ? parts.slice(1).join('_') : idea.file_path.split('/').pop();
      attachmentEl.textContent = filename;

      viewSubmissionBtn.onclick = () => {
        window.open(`/api/submissions/${idea.id}/file`, '_blank');
      };
      viewSubmissionBtn.classList.remove('hidden');
    } else {
      attachmentEl.textContent = 'No attachment';
      viewSubmissionBtn.classList.add('hidden');
    }

    domainBadge.innerHTML = `
      <i class="fas fa-tag text-[10px]"></i>
      ${idea.domain || 'Unspecified'}
    `;

    // Render tags if present
    const tagsContainer = document.getElementById('ideaTagsContainer');
    if (idea.tags && idea.tags.length > 0) {
      tagsContainer.innerHTML = idea.tags.map(tag => `
        <span class="badge-pill bg-purple-50 text-purple-700 border border-purple-200">
          <i class="fas fa-hashtag text-[10px]"></i>
          ${tag}
        </span>
      `).join('');
    } else {
      tagsContainer.innerHTML = '';
    }

    renderStatusBadge(idea.status);

    // Check if we're in incubation mode or if idea is approved
    const isIncubationMode = getMode() === 'incubation' || idea.status === 'approved';
    const reviewPanel = document.getElementById('reviewControlsPanel');
    const incubationPanel = document.getElementById('incubationControlsPanel');

    if (isIncubationMode) {
      // Show incubation controls, hide review controls
      reviewPanel.classList.add('hidden');
      incubationPanel.classList.remove('hidden');

      // Initialize incubation controls with current values
      const stageSelect = document.getElementById('incubationStageSelect');
      const progressRange = document.getElementById('incubationProgressRange');
      const progressValue = document.getElementById('incubationProgressValue');

      stageSelect.value = idea.stage || 'under_incubation';
      progressRange.value = idea.progress_percent || 0;
      progressValue.textContent = `${idea.progress_percent || 0}%`;

      progressRange.oninput = () => {
        progressValue.textContent = `${progressRange.value}%`;
      };
    } else {
      // Show review controls, hide incubation controls
      reviewPanel.classList.remove('hidden');
      incubationPanel.classList.add('hidden');

      // Disable review buttons if already processed
      const reviewButtons = document.getElementById('ideaReviewButtons');
      if (idea.status === 'approved' || idea.status === 'rejected') {
        reviewButtons.querySelectorAll('button').forEach(btn => {
          btn.disabled = true;
          btn.classList.add('opacity-50', 'cursor-not-allowed');
          btn.title = "Decision already made.";
        });
      } else {
        // Wire Review Buttons only if not processed
        reviewButtons.querySelectorAll('[data-review-tag]').forEach(btn => {
          btn.onclick = () => submitDecision(idea.id, btn.dataset.reviewTag);
        });
      }
    }

    ideaContainer.classList.remove('hidden');
  }

  // ---- Feedback Logic ----
  const saveFeedbackBtn = document.getElementById('saveFeedbackBtn');
  const feedbackTextarea = document.getElementById('ideaFeedback');

  saveFeedbackBtn.onclick = async () => {
    const feedback = feedbackTextarea.value.trim();
    if (!feedback) {
      alert('Please enter some feedback.');
      return;
    }

    try {
      const res = await fetch(`/api/faculty/feedback`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': 'Bearer ' + token
        },
        body: JSON.stringify({
          submission_id: submissionId,
          overall_feedback: feedback
        })
      });

      if (!res.ok) throw new Error(await res.text());

      alert('Feedback saved successfully!');
    } catch (e) {
      console.error(e);
      alert('Failed to save feedback: ' + e.message);
    }
  };

  // ---- Incubation Progress Logic ----
  const saveIncubationBtn = document.getElementById('saveIncubationBtn');

  saveIncubationBtn.onclick = async () => {
    const stageSelect = document.getElementById('incubationStageSelect');
    const progressRange = document.getElementById('incubationProgressRange');

    try {
      const res = await fetch(`/api/faculty/incubation/${submissionId}`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': 'Bearer ' + token
        },
        body: JSON.stringify({
          stage: stageSelect.value,
          progress_percent: parseInt(progressRange.value),
          company_id: ''
        })
      });

      if (!res.ok) throw new Error(await res.text());

      alert('Progress updated successfully!');
      window.location.reload();
    } catch (e) {
      console.error(e);
      alert('Failed to update progress: ' + e.message);
    }
  };

  function renderStatusBadge(status) {
    const map = {
      admin_approved: { label: 'Pending Faculty Review', color: 'orange', icon: 'fa-hourglass-half' },
      approved: { label: 'Approved', color: 'green', icon: 'fa-check-circle' },
      rejected: { label: 'Rejected', color: 'red', icon: 'fa-times-circle' }
    };
    const info = map[status] || { label: status, color: 'gray', icon: 'fa-circle' };

    const colorClasses = {
      orange: 'bg-orange-50 text-orange-800 border-orange-200',
      green: 'bg-emerald-50 text-emerald-800 border-emerald-300',
      red: 'bg-rose-50 text-rose-800 border-rose-300',
      gray: 'bg-gray-100 text-gray-700 border-gray-200'
    };

    statusBadge.className = `badge-pill ${colorClasses[info.color]}`;
    statusBadge.innerHTML = `<i class="fas ${info.icon} text-[10px]"></i> ${info.label}`;
  }

  async function submitDecision(id, decision) {
    // Map internal tags to backend statuses
    let apiDecision = decision;
    if (decision === 'accepted') apiDecision = 'approved';
    if (decision === 'needs-improvement') apiDecision = 'needs_improvement';
    // rejected stays as 'rejected'

    try {
      const res = await fetch(`/api/faculty/reviews/${id}/decision`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': 'Bearer ' + token
        },
        body: JSON.stringify({ decision: apiDecision })
      });

      if (!res.ok) throw new Error(await res.text());

      const messages = {
        'approved': 'Submission approved! It has been moved to the Incubation Pipeline.',
        'rejected': 'Submission rejected.',
        'needs_improvement': 'Submission marked as needing improvement. The student will be notified.'
      };
      alert(messages[apiDecision] || `Submission ${apiDecision} successfully!`);
      window.location.reload();
    } catch (e) {
      console.error(e);
      alert('Failed to submit decision: ' + e.message);
    }
  }

  // ---- Navigation ----
  backBtn?.addEventListener('click', () => window.history.back());

  document.getElementById('mobile-menu-btn')
    ?.addEventListener('click', () =>
      document.getElementById('mobile-menu')?.classList.toggle('hidden')
    );

  // ---- Init ----
  const submissionId = getSubmissionId();
  if (!submissionId) {
    ideaNotFound.classList.remove('hidden');
    return;
  }

  const idea = await fetchIdea(submissionId);
  if (!idea || idea.error) {
    ideaNotFound.classList.remove('hidden');
    if (idea && idea.error) {
      ideaNotFound.innerText = `Error: ${idea.status === 404 ? 'Idea not found.' : 'Failed to load idea (' + idea.status + ').'}`;
    }
    return;
  }

  renderIdea(idea);
});