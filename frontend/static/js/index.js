async function loadResources() {
  try {
    const res = await fetch('/api/content/resources/top');
    if (!res.ok) return;

    const data = await res.json();
    const container = document.querySelector('#resources .grid');
    if (!container) return;

    if (!data || data.length === 0) {
      container.innerHTML = '<p class="text-white col-span-3 text-center">No resources available at the moment.</p>';
      return;
    }

    container.innerHTML = data.map((r, i) => `
      <div class="${i % 2 === 0 ? 'bg-white' : 'bg-gray-900 text-white'} p-8 rounded-3xl hover:scale-105 transition shadow-sm hover:shadow-xl duration-300">
        <div class="text-6xl font-black text-orange-primary mb-6">${i + 1}</div>
        <h3 class="text-xl font-bold mb-4 font-clash">${r.title}</h3>
        <p class="${i % 2 === 0 ? 'text-gray-600' : 'text-gray-300'} mb-6 text-sm font-sans">
          ${r.description || 'Access premium innovation materials.'}
        </p>
        <a href="/resources" class="text-orange-primary font-semibold uppercase tracking-widest text-[10px] hover:underline font-sans">
          Learn More
        </a>
      </div>
    `).join('');
  } catch (err) {
    console.error('loadResources failed:', err);
  }
}

async function loadEvents() {
  try {
    const res = await fetch('/api/content/events/upcoming');
    if (!res.ok) return;

    const data = await res.json();
    const list = document.getElementById('eventList');
    if (!list) return;

    if (!data || data.length === 0) {
      list.innerHTML = '<p class="text-white/70">No upcoming events scheduled.</p>';
      return;
    }

    list.innerHTML = data.slice(0, 3).map((e, i) => `
      <div class="flex items-start space-x-4 group">
        <div class="w-10 h-10 bg-white/10 group-hover:bg-white rounded-full flex items-center justify-center text-white group-hover:text-orange-primary font-bold transition-colors">
          ${i + 1}
        </div>
        <div>
          <p class="text-white/90 font-bold font-clash">${e.title}</p>
          <p class="text-white/60 text-xs font-sans">${e.event_date ? new Date(e.event_date).toLocaleDateString() : "TBA"} | ${e.venue || "Global"}</p>
        </div>
      </div>
    `).join('');

    const eventCard = document.getElementById('eventCard');
    if (eventCard && data.length > 0) {
      const e = data[0];
      eventCard.innerHTML = `
        <div class="text-xs font-bold text-orange-primary mb-2 uppercase tracking-[0.2em] font-sans">Featured Event</div>
        <div class="text-2xl font-black text-white mb-4 font-clash">${e.title}</div>
        <div class="space-y-3">
          <div class="flex items-center gap-3 text-gray-400 text-sm">
            <span class="material-symbols-outlined text-sm">calendar_today</span>
            <span>${e.event_date ? new Date(e.event_date).toLocaleDateString() : "Date TBA"}</span>
          </div>
          <div class="flex items-center gap-3 text-gray-400 text-sm">
            <span class="material-symbols-outlined text-sm">location_on</span>
            <span>${e.venue || "Location TBA"}</span>
          </div>
        </div>
        <p class="mt-6 text-gray-500 text-sm leading-relaxed font-sans">${e.description || "Join us for this exciting innovation session."}</p>
      `;
    }
  } catch (err) {
    console.error('loadEvents failed:', err);
  }
}

async function loadShowcaseOverview() {
  try {
    const res = await fetch('/api/submissions/incubation');
    if (!res.ok) return;
    const data = await res.json();

    const counts = {
      under_incubation: 0,
      looking_for_funding: 0,
      found_company: 0
    };

    data.forEach(item => {
      if (counts.hasOwnProperty(item.stage)) {
        counts[item.stage]++;
      }
    });

    const elIncubation = document.getElementById('count-incubation');
    const elFunding = document.getElementById('count-funding');
    const elCompany = document.getElementById('count-company');

    if (elIncubation) elIncubation.textContent = counts.under_incubation;
    if (elFunding) elFunding.textContent = counts.looking_for_funding;
    if (elCompany) elCompany.textContent = counts.found_company;

    // Update hero project counter
    const heroCount = document.getElementById('hero-project-count');
    if (heroCount) {
      const total = data.length;
      heroCount.textContent = total < 10 ? `0 ${total}` : total;
    }

  } catch (err) {
    console.error('loadShowcaseOverview failed:', err);
  }
}

async function loadAboutFeatures() {
  try {
    const res = await fetch('/api/content/about/features');
    if (!res.ok) return;

    const data = await res.json();
    const container = document.getElementById('about-features-container');
    if (!container) return;

    if (!data || data.length === 0) return;

    const icons = ['fa-lightbulb', 'fa-cogs', 'fa-award', 'fa-rocket', 'fa-users'];

    container.innerHTML = data.slice(0, 3).map((f, i) => `
      <div class="flex items-start space-x-4">
        <div class="w-12 h-12 bg-orange-primary rounded-full flex items-center justify-center flex-shrink-0 shadow-lg shadow-orange-primary/20">
          <i class="fas ${icons[i % icons.length]} text-white"></i>
        </div>
        <div>
          <h4 class="text-lg font-bold text-gray-900 mb-2 font-clash uppercase tracking-wider">${f.title}</h4>
          <p class="text-gray-600 font-sans text-sm leading-relaxed">${f.description}</p>
        </div>
      </div>
    `).join('');
  } catch (err) {
    console.error('loadAboutFeatures failed:', err);
  }
}

document.addEventListener('DOMContentLoaded', () => {
  loadResources();
  loadEvents();
  loadShowcaseOverview();
  loadAboutFeatures();
});

