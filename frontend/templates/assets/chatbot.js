/**
 * MAHE Innovation Centre — Chatbot Widget
 * Self-initializing. Exposes window.MaheChatbot with .open(), .close(), .toggle(), .addBotMessage()
 */
(function () {
    'use strict';

    const API_URL = '/api/chat';
    let sessionId = localStorage.getItem('mahe_chat_session') || crypto.randomUUID();
    localStorage.setItem('mahe_chat_session', sessionId);

    /* ---------- DOM creation helpers ---------- */
    function el(tag, attrs, ...children) {
        const e = document.createElement(tag);
        if (attrs) Object.entries(attrs).forEach(([k, v]) => {
            if (k === 'className') e.className = v;
            else if (k.startsWith('on')) e.addEventListener(k.slice(2).toLowerCase(), v);
            else e.setAttribute(k, v);
        });
        children.forEach(c => {
            if (typeof c === 'string') e.appendChild(document.createTextNode(c));
            else if (c) e.appendChild(c);
        });
        return e;
    }

    /* ---------- Build widget DOM ---------- */
    function buildWidget() {
        // Overlay
        const overlay = el('div', { id: 'mahe-chatbot-overlay' });
        overlay.addEventListener('click', () => MaheChatbot.close());

        // Panel
        const panel = el('div', { className: 'mahe-chatbot__panel', id: 'mahe-chatbot-panel' });

        // Header
        const header = el('div', { className: 'mahe-chatbot__header' },
            el('div', { className: 'mahe-chatbot__title' },
                el('div', { className: 'mahe-chatbot__avatar' }, '∞∞'),
                el('div', { className: 'mahe-chatbot__titleText' },
                    el('strong', null, 'MAHE Assistant'),
                    el('span', null, 'Ask anything about MAHE Innovation Centre')
                )
            ),
            el('div', { className: 'mahe-chatbot__headerActions' },
                el('button', { className: 'mahe-chatbot__iconBtn', id: 'mahe-chatbot-clear', 'aria-label': 'Clear chat' },
                    el('i', { className: 'fas fa-trash-alt' })
                ),
                el('button', { className: 'mahe-chatbot__iconBtn', id: 'mahe-chatbot-close', 'aria-label': 'Close chat' },
                    el('i', { className: 'fas fa-times' })
                )
            )
        );

        // Messages container
        const messages = el('div', { className: 'mahe-chatbot__messages', id: 'mahe-chatbot-messages' });

        // Composer
        const fileInput = el('input', { type: 'file', className: 'mahe-chatbot__fileInput', id: 'mahe-chatbot-file-input', multiple: '', accept: '*/*' });
        const filePreview = el('div', { id: 'mahe-chatbot-file-preview' });
        const textarea = el('textarea', {
            className: 'mahe-chatbot__input',
            id: 'mahe-chatbot-input',
            rows: '1',
            placeholder: 'Type a message… (Enter to send)'
        });
        const sendBtn = el('button', { className: 'mahe-chatbot__send', id: 'mahe-chatbot-send', type: 'button', 'aria-label': 'Send' },
            el('i', { className: 'fas fa-paper-plane' })
        );

        const inputWrap = el('div', { className: 'mahe-chatbot__inputWrap' },
            textarea,
            sendBtn
        );

        const composer = el('div', { className: 'mahe-chatbot__composer' },
            filePreview,
            inputWrap
        );

        panel.appendChild(header);
        panel.appendChild(messages);
        panel.appendChild(composer);

        // Launcher FAB
        const launcher = el('button', { id: 'mahe-chatbot-launcher', 'aria-label': 'Open Chatbot' });
        launcher.innerHTML = '<svg viewBox="0 0 24 24"><path d="M20 2H4c-1.1 0-2 .9-2 2v18l4-4h14c1.1 0 2-.9 2-2V4c0-1.1-.9-2-2-2zm0 14H5.2L4 17.2V4h16v12z"/><path d="M7 9h10v2H7zm0-3h10v2H7z"/></svg>';

        document.body.appendChild(overlay);
        document.body.appendChild(panel);
        document.body.appendChild(launcher);

        return { overlay, panel, messages, textarea, sendBtn, launcher };
    }

    /* ---------- Render helpers ---------- */
    function parseButtons(text) {
        // Convert [BUTTON:label|url] to HTML links
        return text.replace(/\[BUTTON:([^|]+)\|([^\]]+)\]/g,
            '<a class="mahe-chatbot__link-btn" href="$2">$1</a>');
    }

    function addMessage(container, role, text) {
        const div = el('div', {
            className: 'mahe-chatbot__msg mahe-chatbot__msg--' + (role === 'user' ? 'user' : 'bot')
        });
        if (role === 'bot') {
            div.innerHTML = parseButtons(text.replace(/\n/g, '<br>'));
        } else {
            div.textContent = text;
        }
        container.appendChild(div);
        container.scrollTop = container.scrollHeight;
    }

    function showTyping(container) {
        const typing = el('div', { className: 'mahe-chatbot__typing', id: 'mahe-chatbot-typing' },
            el('span'), el('span'), el('span')
        );
        container.appendChild(typing);
        container.scrollTop = container.scrollHeight;
        return typing;
    }

    function removeTyping() {
        const t = document.getElementById('mahe-chatbot-typing');
        if (t) t.remove();
    }

    /* ---------- API call ---------- */
    async function sendToAPI(message) {
        try {
            const res = await fetch(API_URL, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                    'X-Session-ID': sessionId
                },
                body: JSON.stringify({ message, session_id: sessionId })
            });
            if (!res.ok) throw new Error('Server error ' + res.status);
            const data = await res.json();
            if (data.session_id) {
                sessionId = data.session_id;
                localStorage.setItem('mahe_chat_session', sessionId);
            }
            return data.response || 'Sorry, I could not process your request.';
        } catch (err) {
            console.error('Chatbot API error:', err);
            return "I'm having trouble connecting right now. Please try again in a moment.";
        }
    }

    /* ---------- Init ---------- */
    function init() {
        // Don't double-init
        if (window.MaheChatbot && window.MaheChatbot._initialized) return;

        const { overlay, panel, messages, textarea, sendBtn, launcher } = buildWidget();
        let isOpen = false;

        function openChat() {
            panel.classList.add('open');
            overlay.classList.add('visible');
            launcher.style.display = 'none';
            isOpen = true;
            textarea.focus();

            // Welcome message on first open
            if (messages.children.length === 0) {
                addMessage(messages, 'bot',
                    "Hi! 👋 I'm the MAHE Innovation Centre Assistant. Ask me about our events, resources, programs, or anything about MiC!");
            }
        }

        function closeChat() {
            panel.classList.remove('open');
            overlay.classList.remove('visible');
            launcher.style.display = 'flex';
            isOpen = false;
        }

        async function handleSend() {
            const text = textarea.value.trim();
            if (!text) return;

            addMessage(messages, 'user', text);
            textarea.value = '';
            textarea.style.height = 'auto';

            const typingEl = showTyping(messages);

            const reply = await sendToAPI(text);

            removeTyping();
            addMessage(messages, 'bot', reply);
        }

        // Event bindings
        launcher.addEventListener('click', openChat);
        document.getElementById('mahe-chatbot-close').addEventListener('click', closeChat);
        document.getElementById('mahe-chatbot-clear').addEventListener('click', () => {
            messages.innerHTML = '';
            sessionId = crypto.randomUUID();
            localStorage.setItem('mahe_chat_session', sessionId);
            addMessage(messages, 'bot',
                "Chat cleared! 🧹 How can I help you today?");
        });

        sendBtn.addEventListener('click', handleSend);
        textarea.addEventListener('keydown', (e) => {
            if (e.key === 'Enter' && !e.shiftKey) {
                e.preventDefault();
                handleSend();
            }
        });

        // Auto-resize textarea
        textarea.addEventListener('input', () => {
            textarea.style.height = 'auto';
            textarea.style.height = Math.min(textarea.scrollHeight, 100) + 'px';
        });

        // Expose global API
        window.MaheChatbot = {
            _initialized: true,
            open: openChat,
            close: closeChat,
            toggle: () => isOpen ? closeChat() : openChat(),
            addBotMessage: (text) => addMessage(messages, 'bot', text)
        };
    }

    // Wait for DOM then initialise
    if (document.readyState === 'loading') {
        document.addEventListener('DOMContentLoaded', init);
    } else {
        init();
    }
})();
