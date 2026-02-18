(function () {
  const STORAGE_KEY = 'mic_faculty_state_v1';

  const defaultState = {
    faculty: {
      id: 'fac-001',
      name: 'Arun Y Patil',
      email: 'arun.patil@mic.edu',
      department: 'Department of Computer Science & Engineering',
      canManageCommittee: true
    },
    ideas: [
      {
        id: 'idea-1',
        title: 'AI-driven Learning Analytics Platform',
        student: 'Team Alpha (B.Tech CSE)',
        status: 'pending', // pending | approved | rejected
        requiresReview: true,
        submittedOn: 'March 15, 2025',
        domain: 'Artificial Intelligence',
        description: 'A platform that uses ML models to analyze learning patterns and provide insights to faculty and students.',
        attachmentName: 'ai_learning_analytics.pdf',
        attachmentUrl: '',
        incubationStage: 'Faculty Review',
        progressPercent: 45
      },
      {
        id: 'idea-2',
        title: 'Smart Campus Energy Optimizer',
        student: 'Rahul & Meera (EEE + CSE)',
        status: 'approved',
        requiresReview: false,
        submittedOn: 'March 10, 2025',
        domain: 'IoT / Sustainability',
        description: 'IoT-based solution to monitor and optimize power usage across the campus buildings.',
        attachmentName: 'smart_campus_energy.pdf',
        attachmentUrl: '',
        incubationStage: 'Incubation - Prototype',
        progressPercent: 72
      },
      {
        id: 'idea-3',
        title: 'Assistive AR Navigation for Visually Impaired',
        student: 'Innovation Club XR',
        status: 'pending',
        requiresReview: true,
        submittedOn: 'March 20, 2025',
        domain: 'AR / Accessibility',
        description: 'An AR navigation assistant that provides audio and haptic feedback to help visually impaired users navigate indoor spaces.',
        attachmentName: 'assistive_ar_navigation.pdf',
        attachmentUrl: '',
        incubationStage: 'Initial Screening',
        progressPercent: 30
      },
      {
        id: 'idea-4',
        title: 'Blockchain-based Transcript Verification',
        student: 'Final Year B.Tech IT',
        status: 'rejected',
        requiresReview: false,
        submittedOn: 'February 22, 2025',
        domain: 'Blockchain / Records',
        description: 'Decentralized system to verify academic transcripts securely using blockchain.',
        attachmentName: 'transcript_blockchain.pdf',
        attachmentUrl: '',
        incubationStage: 'Rejected',
        progressPercent: 15
      }
    ],
    committeeMembers: [
      { id: 'mem-1', name: 'Arun Y Patil', role: 'Faculty Coordinator', email: 'arun.patil@mic.edu' },
      { id: 'mem-2', name: 'Dr. Neha Sharma', role: 'Co-Coordinator (Research)', email: 'neha.sharma@mic.edu' },
      { id: 'mem-3', name: 'Prof. Kunal Rao', role: 'Industry Liaison', email: 'kunal.rao@mic.edu' },
      { id: 'mem-4', name: 'Dr. Priya Menon', role: 'Innovation Mentor', email: 'priya.menon@mic.edu' }
    ],
    events: [
      {
        id: 'evt-1',
        title: 'MIC Ideathon 2025 - Final Jury',
        role: 'Jury Member',
        date: 'April 5, 2025',
        rsvpBy: 'March 28, 2025',
        status: 'pending', // pending | accepted | declined
        location: 'Innovation Centre Auditorium'
      },
      {
        id: 'evt-2',
        title: 'Startup Clinic - AI & Data Track',
        role: 'Mentor',
        date: 'April 18, 2025',
        rsvpBy: 'April 10, 2025',
        status: 'accepted',
        location: 'MIC Co-working Lab'
      },
      {
        id: 'evt-3',
        title: 'Engineering Design Expo - Judge',
        role: 'Judge',
        date: 'May 2, 2025',
        rsvpBy: 'April 20, 2025',
        status: 'pending',
        location: 'New Academic Block, Hall 3'
      }
    ]
  };

  function loadState() {
    try {
      const raw = localStorage.getItem(STORAGE_KEY);
      if (!raw) return structuredClone(defaultState);
      const parsed = JSON.parse(raw);
      // Shallow merge to avoid breaking if schema changes
      return Object.assign({}, defaultState, parsed);
    } catch (e) {
      console.warn('Failed to load faculty state, using defaults', e);
      return structuredClone(defaultState);
    }
  }

  function saveState(state) {
    try {
      localStorage.setItem(STORAGE_KEY, JSON.stringify(state));
    } catch (e) {
      console.warn('Failed to save faculty state', e);
    }
  }

  function computeIdeaStats(state) {
    const total = state.ideas.length;
    let pending = 0, approved = 0, rejected = 0, requiresReview = 0;
    state.ideas.forEach(i => {
      if (i.status === 'pending') pending++;
      if (i.status === 'approved') approved++;
      if (i.status === 'rejected') rejected++;
      if (i.status === 'pending' && i.requiresReview) requiresReview++;
    });
    return { total, pending, approved, rejected, requiresReview };
  }

  function computeEventStats(state) {
    let pending = 0, accepted = 0, declined = 0;
    state.events.forEach(e => {
      if (e.status === 'pending') pending++;
      if (e.status === 'accepted') accepted++;
      if (e.status === 'declined') declined++;
    });
    return { pending, accepted, declined };
  }

  const FacultyStore = {
    getState() {
      if (!this._state) {
        this._state = loadState();
      }
      return this._state;
    },

    resetState() {
      this._state = structuredClone(defaultState);
      saveState(this._state);
    },

    updateFaculty(partial) {
      const state = this.getState();
      state.faculty = Object.assign({}, state.faculty, partial);
      saveState(state);
      return state.faculty;
    },

    setCommitteeMembers(members) {
      const state = this.getState();
      state.committeeMembers = members.map((m, idx) => ({
        id: m.id || `mem-${idx + 1}`,
        name: m.name,
        role: m.role,
        email: m.email
      }));
      saveState(state);
      return state.committeeMembers;
    },

    addCommitteeMember(member) {
      const state = this.getState();
      const next = {
        id: member.id || `mem-${Date.now()}`,
        name: member.name,
        role: member.role,
        email: member.email
      };
      state.committeeMembers.push(next);
      saveState(state);
      return next;
    },

    removeCommitteeMember(id) {
      const state = this.getState();
      state.committeeMembers = state.committeeMembers.filter(m => m.id !== id);
      saveState(state);
      return state.committeeMembers;
    },

    upsertIdea(idea) {
      const state = this.getState();
      const idx = state.ideas.findIndex(i => i.id === idea.id);
      if (idx === -1) {
        const id = idea.id || `idea-${Date.now()}`;
        const next = Object.assign({ id }, idea);
        state.ideas.push(next);
        saveState(state);
        return next;
      } else {
        state.ideas[idx] = Object.assign({}, state.ideas[idx], idea);
        saveState(state);
        return state.ideas[idx];
      }
    },

    updateIdeaStatus(id, status, options) {
      const state = this.getState();
      const idea = state.ideas.find(i => i.id === id);
      if (!idea) return null;
      idea.status = status;
      if (options && typeof options.requiresReview === 'boolean') {
        idea.requiresReview = options.requiresReview;
      }
      saveState(state);
      return idea;
    },

    upsertEvent(event) {
      const state = this.getState();
      const idx = state.events.findIndex(e => e.id === event.id);
      if (idx === -1) {
        const id = event.id || `evt-${Date.now()}`;
        const next = Object.assign({ id }, event);
        state.events.push(next);
        saveState(state);
        return next;
      } else {
        state.events[idx] = Object.assign({}, state.events[idx], event);
        saveState(state);
        return state.events[idx];
      }
    },

    updateEventStatus(id, status) {
      const state = this.getState();
      const event = state.events.find(e => e.id === id);
      if (!event) return null;
      event.status = status;
      saveState(state);
      return event;
    },

    getIdeaStats() {
      return computeIdeaStats(this.getState());
    },

    getEventStats() {
      return computeEventStats(this.getState());
    }
  };

  window.FacultyStore = FacultyStore;
})();
