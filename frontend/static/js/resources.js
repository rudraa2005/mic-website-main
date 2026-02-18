/* ================= Fetch helpers ================= */

async function fetchTopResources() {
  try {
    const res = await fetch('/api/content/resources/top');
    if (!res.ok) {
      console.error('Failed to load top resources');
      return [];
    }
    return await res.json();
  } catch (err) {
    console.error('Top resources fetch error:', err);
    return [];
  }
}

async function fetchAllResources() {
  try {
    const res = await fetch('/api/content/resources');
    if (!res.ok) {
      console.error('Failed to load resources');
      return [];
    }
    return await res.json();
  } catch (err) {
    console.error('Resources fetch error:', err);
    return [];
  }
}

/* ================= Featured Resources ================= */

function renderFeaturedResources(items) {
  const container = document.getElementById('featuredResources');
  if (!container) return;

  if (!items.length) {
    container.innerHTML =
      `<p class="text-sm text-gray-400">No featured resources available</p>`;
    return;
  }

  container.innerHTML = items.map(r => `
    <div class="bg-black text-white rounded-2xl p-6 shadow-xl">
      <h3 class="text-lg font-bold mb-2">
        ${r.title}
      </h3>

      <p class="text-sm text-gray-300 mb-4">
        ${r.description || 'No description available'}
      </p>

      ${
        r.file_url
          ? `<a href="${r.file_url}" target="_blank"
               class="inline-block text-xs font-semibold text-orange-primary">
               VIEW RESOURCE →
             </a>`
          : ''
      }
    </div>
  `).join('');
}

/* ================= Resource Grid ================= */

function renderResourceGrid(items) {
  const container = document.getElementById('resourceGrid');
  if (!container) return;

  if (!items.length) {
    container.innerHTML =
      `<p class="text-sm text-gray-500">No resources found</p>`;
    return;
  }

  container.innerHTML = items.map((r, i) => `
    <div class="rounded-xl p-6 ${i % 2 === 0 ? 'bg-gray-200' : 'bg-black text-white'}">
      <div class="flex items-center gap-3 mb-3">
        <div class="w-7 h-7 ${i % 2 === 0 ? 'bg-gray-800' : 'bg-gray-600'} rounded"></div>
        <span class="text-xs font-bold">
          ${String(i + 1).padStart(2, '0')}
        </span>
      </div>

      <p class="text-sm mb-4">
        ${r.description || ''}
      </p>

      ${
        r.file_url
          ? `<a href="${r.file_url}" target="_blank"
               class="text-xs font-semibold text-orange-primary">
               OPEN →
             </a>`
          : ''
      }
    </div>
  `).join('');
}

/* ================= Page Boot ================= */

document.addEventListener('DOMContentLoaded', async () => {
  const [featured, resources] = await Promise.all([
    fetchTopResources(),
    fetchAllResources()
  ]);

  renderFeaturedResources(featured);
  renderResourceGrid(resources);
});