document.addEventListener('DOMContentLoaded', async () => {
    const portfolioGrid = document.getElementById('portfolioGrid');
    const emptyState = document.getElementById('emptyState');
    const updateModal = document.getElementById('updateModal');
    const closeModal = document.getElementById('closeModal');
    const saveUpdate = document.getElementById('saveUpdate');

    const stageSelect = document.getElementById('stageSelect');
    const progressRange = document.getElementById('progressRange');
    const progressValue = document.getElementById('progressValue');
    const companySelect = document.getElementById('companySelect');

    let currentSubmissionId = null;
    const token = localStorage.getItem('authToken');

    if (!token) {
        window.location.href = '/login.html';
        return;
    }

    progressRange.oninput = () => {
        progressValue.textContent = `${progressRange.value}%`;
    };

    async function fetchData() {
        try {
            const [portfolioRes, companiesRes] = await Promise.all([
                fetch('/api/faculty/incubation', { headers: { 'Authorization': 'Bearer ' + token } }),
                fetch('/api/faculty/companies', { headers: { 'Authorization': 'Bearer ' + token } })
            ]);

            const portfolio = await portfolioRes.json();
            const companies = await companiesRes.json();

            renderPortfolio(portfolio);
            renderCompanies(companies);
        } catch (err) {
            console.error(err);
        }
    }

    function renderPortfolio(items) {
        const incubationGrid = document.getElementById('incubationGrid');
        const fundingGrid = document.getElementById('fundingGrid');
        const fundedGrid = document.getElementById('fundedGrid');
        const incubationEmpty = document.getElementById('incubationEmpty');
        const fundingEmpty = document.getElementById('fundingEmpty');
        const fundedEmpty = document.getElementById('fundedEmpty');
        const incubationCount = document.getElementById('incubationCount');
        const fundingCount = document.getElementById('fundingCount');
        const fundedCount = document.getElementById('fundedCount');

        if (!items || items.length === 0) {
            emptyState.classList.remove('hidden');
            incubationEmpty.classList.add('hidden');
            fundingEmpty.classList.add('hidden');
            fundedEmpty.classList.add('hidden');
            return;
        }
        emptyState.classList.add('hidden');

        // Separate items by stage
        const incubationItems = items.filter(i => i.stage === 'under_incubation' || !i.stage);
        const fundingItems = items.filter(i => i.stage === 'looking_for_funding');
        const fundedItems = items.filter(i => i.stage === 'found_company');

        // Update counts
        incubationCount.textContent = incubationItems.length;
        fundingCount.textContent = fundingItems.length;
        fundedCount.textContent = fundedItems.length;

        const renderCard = (item, colorClass) => `
            <div class="glass-card rounded-[2rem] p-8 hover:shadow-xl transition-all border border-slate-100 group">
                <div class="flex justify-between items-start mb-6">
                    <div>
                        <h4 class="text-xl font-black text-slate-900 uppercase tracking-tight mb-1">${item.title}</h4>
                        <p class="text-xs font-bold ${colorClass} uppercase tracking-widest">${getStageLabel(item.stage)}</p>
                    </div>
                    <div class="flex gap-2">
                        <a href="faculty-idea.html?id=${item.submission_id}&mode=incubation" 
                            class="w-10 h-10 rounded-full bg-blue-100 flex items-center justify-center hover:bg-blue-600 hover:text-white text-blue-600 transition-all"
                            title="View Details">
                            <i class="fas fa-eye text-xs"></i>
                        </a>
                        <button onclick="openUpdateModal('${item.submission_id}', '${item.title}', '${item.stage || 'under_incubation'}', ${item.progress_percent || 0})" 
                            class="w-10 h-10 rounded-full bg-slate-100 flex items-center justify-center hover:bg-orange-600 hover:text-white transition-all"
                            title="Update Progress">
                            <i class="fas fa-edit text-xs"></i>
                        </button>
                    </div>
                </div>
                
                <div class="space-y-4 mb-8">
                    <div class="flex justify-between text-[10px] font-bold uppercase border-b border-slate-50 pb-2">
                        <span class="text-slate-400">Current Progress</span>
                        <span class="text-slate-900">${item.progress_percent || 0}%</span>
                    </div>
                    <div class="w-full bg-slate-100 h-2 rounded-full overflow-hidden">
                        <div class="bg-orange-500 h-full transition-all duration-500" style="width: ${item.progress_percent || 0}%"></div>
                    </div>
                </div>

                <div class="flex items-center gap-3">
                    <div class="w-8 h-8 rounded-full bg-slate-100 flex items-center justify-center text-slate-400">
                      <i class="fas fa-user text-[10px]"></i>
                    </div>
                    <span class="text-xs font-medium text-slate-600">${item.student}</span>
                </div>
            </div>
        `;

        const getStageLabel = (stage) => {
            const labels = {
                'under_incubation': 'Under Incubation',
                'looking_for_funding': 'Looking for Funding',
                'found_company': 'Funded Company'
            };
            return labels[stage] || 'Under Incubation';
        };

        // Render each section
        incubationGrid.innerHTML = incubationItems.map(i => renderCard(i, 'text-orange-600')).join('');
        fundingGrid.innerHTML = fundingItems.map(i => renderCard(i, 'text-blue-600')).join('');
        fundedGrid.innerHTML = fundedItems.map(i => renderCard(i, 'text-emerald-600')).join('');

        // Show/hide empty states
        incubationEmpty.classList.toggle('hidden', incubationItems.length > 0);
        fundingEmpty.classList.toggle('hidden', fundingItems.length > 0);
        fundedEmpty.classList.toggle('hidden', fundedItems.length > 0);
    }

    function renderCompanies(companies) {
        companySelect.innerHTML = '<option value="">No Affiliation</option>' +
            companies.map(c => `<option value="${c.id}">${c.name}</option>`).join('');
    }

    window.openUpdateModal = (id, title, stage, progress) => {
        currentSubmissionId = id;
        document.getElementById('modalProjectTitle').textContent = title;
        stageSelect.value = stage;
        progressRange.value = progress;
        progressValue.textContent = `${progress}%`;
        updateModal.classList.remove('hidden');
    };

    closeModal.onclick = () => updateModal.classList.add('hidden');

    saveUpdate.onclick = async () => {
        try {
            const res = await fetch(`/api/faculty/incubation/${currentSubmissionId}`, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                    'Authorization': 'Bearer ' + token
                },
                body: JSON.stringify({
                    stage: stageSelect.value,
                    progress_percent: parseInt(progressRange.value),
                    company_id: companySelect.value
                })
            });

            if (!res.ok) throw new Error('Update failed');

            updateModal.classList.add('hidden');
            fetchData();
        } catch (err) {
            alert(err.message);
        }
    };

    fetchData();
});
