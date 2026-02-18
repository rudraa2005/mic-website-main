document.addEventListener('DOMContentLoaded', async () => {
    const API_URL = '/api/submissions/incubation';

    // Select containers
    const stage1Container = document.querySelector('#stage1-projects');
    const stage2Container = document.querySelector('#stage2-projects');
    const stage3Container = document.querySelector('#stage3-projects');
    const featuredContainer = document.querySelector('#featured-project-container');

    async function fetchPipeline() {
        try {
            const res = await fetch(API_URL);
            if (!res.ok) throw new Error('Failed to fetch');
            const data = await res.json();
            renderPipeline(data);
        } catch (err) {
            console.error(err);
        }
    }

    function renderPipeline(submissions) {
        // Clear containers if they exist
        if (stage1Container) stage1Container.innerHTML = '';
        if (stage2Container) stage2Container.innerHTML = '';
        if (stage3Container) stage3Container.innerHTML = '';
        if (featuredContainer) featuredContainer.innerHTML = '';
        const partnersContainer = document.querySelector('#partners-container');
        if (partnersContainer) partnersContainer.innerHTML = '';

        if (!submissions || submissions.length === 0) {
            // Hide all wrappers if no data
            ['stage1-wrapper', 'stage2-wrapper', 'stage3-wrapper'].forEach(id => {
                const el = document.getElementById(id);
                if (el) el.style.display = 'none';
            });
            return;
        }

        // Render Hero Stats
        const heroStats = document.querySelector('#project-count-hero');
        if (heroStats) heroStats.textContent = submissions.length.toString().padStart(2, '0');

        // Render Featured Project (take the first one)
        renderFeatured(submissions[0]);

        // Track unique partners and project counts for visibility
        const partners = new Map();
        const counts = { stage1: 0, stage2: 0, stage3: 0 };

        // Render all projects into stages
        submissions.forEach(submission => {
            const html = createProjectCard(submission);

            if (submission.company_name) {
                partners.set(submission.company_name, submission.company_logo);
            }

            // Map stage to container
            if (submission.stage === 'under_incubation') {
                if (stage1Container) {
                    stage1Container.insertAdjacentHTML('beforeend', html);
                    counts.stage1++;
                }
            } else if (submission.stage === 'looking_for_funding') {
                if (stage2Container) {
                    stage2Container.insertAdjacentHTML('beforeend', html);
                    counts.stage2++;
                }
            } else if (submission.stage === 'found_company') {
                if (stage3Container) {
                    stage3Container.insertAdjacentHTML('beforeend', html);
                    counts.stage3++;
                }
            } else {
                if (stage1Container) {
                    stage1Container.insertAdjacentHTML('beforeend', html);
                    counts.stage1++;
                }
            }
        });

        // Toggle Wrapper Visibility
        const s1Wrapper = document.getElementById('stage1-wrapper');
        const s2Wrapper = document.getElementById('stage2-wrapper');
        const s3Wrapper = document.getElementById('stage3-wrapper');

        if (s1Wrapper) s1Wrapper.style.display = counts.stage1 > 0 ? 'block' : 'none';
        if (s2Wrapper) s2Wrapper.style.display = counts.stage2 > 0 ? 'block' : 'none';
        if (s3Wrapper) s3Wrapper.style.display = counts.stage3 > 0 ? 'block' : 'none';

        // Render Partners
        if (partnersContainer && partners.size > 0) {
            partners.forEach((logo, name) => {
                const partnerHtml = `
                    <div class="flex flex-col items-center gap-3">
                        ${logo ? `<img src="${logo}" class="w-12 h-12 object-contain grayscale opacity-80 group-hover:opacity-100 transition-opacity">` : '<span class="material-symbols-outlined text-4xl">corporate_fare</span>'}
                        <span class="text-[10px] font-bold tracking-widest font-sans uppercase">${name}</span>
                    </div>
                `;
                partnersContainer.insertAdjacentHTML('beforeend', partnerHtml);
            });
        }
    }

    function renderFeatured(s) {
        if (!featuredContainer) return;

        const featuredHtml = `
            <div class="relative w-full aspect-square max-w-lg ml-auto">
                <div class="absolute right-0 top-0 w-4/5 h-full bg-navy-dark rounded-t-[100px] rounded-b-lg shadow-2xl overflow-hidden group border-8 border-white">
                    <img alt="${s.title}"
                        class="w-full h-full object-cover opacity-60 group-hover:scale-110 transition-transform duration-700"
                        src="${s.file_path || 'https://images.unsplash.com/photo-1485827404703-89b55fcc595e?auto=format&fit=crop&w=800&q=80'}" />
                    <div class="absolute bottom-10 left-10 right-10 text-white">
                        <span class="text-primary text-xs font-bold uppercase tracking-widest mb-3 block font-sans">Now Trending</span>
                        <h3 class="text-3xl font-bold font-display">${s.title}</h3>
                        <p class="text-base text-slate-300 mt-2 font-sans">${s.description || 'Innovation in progress'}</p>
                    </div>
                </div>
                <div class="absolute -bottom-6 left-20 w-48 h-32 bg-primary/20 rounded-xl -z-10 blur-2xl"></div>
                <div class="absolute top-1/4 -right-6 w-16 h-16 bg-primary rounded-full shadow-2xl flex items-center justify-center">
                    <span class="material-symbols-outlined text-white">bolt</span>
                </div>
                <div class="absolute top-0 right-0 p-4 text-slate-400 font-sans font-medium tracking-tighter text-vertical text-[10px] uppercase">
                    Innovation Hub
                </div>
            </div>
        `;
        featuredContainer.innerHTML = featuredHtml;
    }

    function createProjectCard(s) {
        const companyHtml = s.company_name ? `
            <div class="mt-6 pt-6 border-t border-slate-100 flex items-center gap-3">
                ${s.company_logo ? `<img src="${s.company_logo}" class="w-6 h-6 object-contain grayscale opacity-50">` : '<i class="fas fa-building text-slate-300"></i>'}
                <span class="text-[10px] font-bold text-slate-400 uppercase tracking-widest">Funded by ${s.company_name}</span>
            </div>
        ` : '';

        return `
            <div class="group bg-slate-50 p-10 rounded-3xl border border-slate-200 hover:border-primary transition-all duration-500 shadow-sm hover:shadow-2xl hover:-translate-y-2">
                <div class="w-12 h-12 bg-white rounded-xl flex items-center justify-center mb-8 shadow-sm">
                    <span class="material-symbols-outlined text-primary text-2xl">${s.stage === 'found_company' ? 'rocket_launch' : (s.stage === 'looking_for_funding' ? 'trending_up' : 'lightbulb')}</span>
                </div>
                <h4 class="text-xl font-bold mb-4 font-display uppercase tracking-tight">${s.title}</h4>
                <p class="text-slate-600 text-sm mb-8 leading-relaxed font-sans">${s.description || 'Innovation in progress.'}</p>
                <div class="flex flex-wrap gap-2">
                    <span class="bg-primary/10 text-primary text-[9px] px-3 py-1 rounded-full font-bold uppercase tracking-wider font-sans">${s.stage || 'Incubating'}</span>
                </div>
                ${companyHtml}
            </div>
        `;
    }

    fetchPipeline();
});
