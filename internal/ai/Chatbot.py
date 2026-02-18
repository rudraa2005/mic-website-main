from groq import Groq
import datetime
from dotenv import dotenv_values
import time
import os
import json
import uuid

# Load env — try multiple paths for .env
env_vars = {}
for p in ["../../.env", ".env", "../../../.env"]:
    loaded = dotenv_values(p)
    if loaded:
        env_vars = loaded
        break

# Also check os.environ as fallback
Assistantname = env_vars.get("Assistantname", os.environ.get("Assistantname", "MAHE Innovation Centre Assistant"))
GroqAPIKey = env_vars.get("GroqAPIKey", env_vars.get("GROQ_API_KEY", os.environ.get("GROQ_API_KEY", "")))

client = None
if GroqAPIKey and GroqAPIKey != "your-groq-api-key-here":
    try:
        client = Groq(api_key=GroqAPIKey)
        print(f"[CHATBOT] Groq client initialized successfully")
    except Exception as e:
        print(f"[CHATBOT] Failed to initialize Groq client: {e}")
        client = None
else:
    print(f"[CHATBOT] No valid Groq API key found, using fallback mode")

# In-memory session store (simple dict keyed by session_id)
_sessions = {}

System = f"""You are {Assistantname}, the official AI assistant for MAHE Innovation Centre (MiC).

Your role is to help visitors and users with information about:
- MAHE Innovation Centre (MiC) - Manipal's Innovation Centre
- Events, workshops, and programs organized by MiC
- Resources, toolkits, guides, and mentorship programs
- Contact information and how to get involved
- Innovation, creation, and incubation programs
- MAHE SID and SCHAP e-Cell initiatives

Key guidelines:
- ONLY answer questions related to MAHE Innovation Centre and its website content
- If asked about topics unrelated to MiC, politely redirect: "I'm here to help with questions about MAHE Innovation Centre. Please ask me about our events, resources, programs, or how to get involved with MiC."
- Respond only in English, even if questions are in other languages
- Keep responses DIRECT and CONCISE - avoid lengthy explanations
- Do NOT use asterisks (*) for formatting
- Maintain a warm, friendly, and professional tone
- Keep responses under 3-4 sentences when possible

About MAHE Innovation Centre (MiC):
- MiC stands for MAHE Innovation Centre, located at Manipal Academy of Higher Education (MAHE)
- It's Manipal's premier hub for innovation, entrepreneurship, and interdisciplinary collaboration
- Provides financial aid and funding opportunities to aspiring entrepreneurs
- Offers incubation programs through MAHE SID (Society for Innovation and Development)
- Provides mentorship and guidance through SCHAP e-Cell initiatives
- Organizes events, workshops, hackathons, and provides resources for innovators
- Focuses on fostering creativity, supporting startups, and building an innovation ecosystem

LINK PROVISION: When appropriate, suggest relevant pages in this exact format:
- For events: [BUTTON:Events Page|/events]
- For resources: [BUTTON:Resources Page|/resources]
- For general info: [BUTTON:About Page|/about]
- For contact: [BUTTON:Contact Page|/contact]
- For home: [BUTTON:Home Page|/]"""

SystemChatBot = [
    {"role": "system", "content": System}
]


def _get_session(session_id):
    """Get or create in-memory session"""
    if session_id not in _sessions:
        _sessions[session_id] = {
            "history": [],
            "context": {
                "current_topic": None,
                "last_question_type": None,
            }
        }
    return _sessions[session_id]


def is_website_related(query):
    """Check if the query is related to MAHE Innovation Centre website"""
    query_lower = query.lower()
    mic_keywords = [
        'mahe', 'mic', 'innovation centre', 'innovation center', 'manipal',
        'event', 'events', 'workshop', 'workshops', 'program', 'programs',
        'resource', 'resources', 'toolkit', 'toolkits', 'guide', 'guides',
        'mentorship', 'incubation', 'incubator', 'entrepreneur', 'entrepreneurship',
        'sid', 'schap', 'e-cell', 'ecell', 'contact', 'about', 'team',
        'funding', 'financial aid', 'startup', 'startups', 'collaboration',
        'what is', 'what exactly is', 'what does', 'explain', 'tell me about',
        'hello', 'hi', 'hey', 'help', 'how', 'who', 'when', 'where',
    ]
    return any(keyword in query_lower for keyword in mic_keywords)


def get_realtime_information():
    """Get current date and time information"""
    now = datetime.datetime.now()
    return f"Current date/time: {now.strftime('%A, %d %B %Y %H:%M')}"


def get_fallback_response(query):
    """Provide fallback responses when API is not available"""
    query_lower = query.lower()

    if any(phrase in query_lower for phrase in ['what is mic', 'what exactly is mic', 'what does mic stand for', 'what is mahe innovation centre']):
        return "MiC stands for MAHE Innovation Centre, Manipal's premier hub for innovation and entrepreneurship. We provide funding, incubation programs, and mentorship to aspiring entrepreneurs. [BUTTON:About Page|/about]"

    elif any(word in query_lower for word in ['event', 'events', 'workshop', 'program']):
        return "We host various events including workshops, hackathons, and innovation showcases. [BUTTON:Events Page|/events]"

    elif any(word in query_lower for word in ['resource', 'resources', 'toolkit', 'guide']):
        return "We provide numerous resources for innovators and entrepreneurs including toolkits, guides, and mentorship materials. [BUTTON:Resources Page|/resources]"

    elif any(word in query_lower for word in ['contact', 'reach', 'get in touch']):
        return "You can contact us through our Contact page or reach out via email. We're here to help! [BUTTON:Contact Page|/contact]"

    elif any(word in query_lower for word in ['about', 'who we are']):
        return "MAHE Innovation Centre is Manipal's hub for innovation and entrepreneurship. [BUTTON:About Page|/about]"

    elif any(word in query_lower for word in ['incubation', 'startup', 'funding']):
        return "We offer incubation support through MAHE SID and provide financial aid to entrepreneurs. [BUTTON:Resources Page|/resources]"

    elif any(word in query_lower for word in ['hello', 'hi', 'hey']):
        return "Hello! Welcome to MAHE Innovation Centre. How can I help you today? You can ask about our events, resources, programs, or anything else about MiC!"

    else:
        return "I'm here to help with questions about MAHE Innovation Centre. Please ask me about our events, resources, programs, or how to get involved with MiC."


def clean_response(response):
    """Remove repetitive sentences and formatting artefacts"""
    if not response:
        return response
    response = response.replace("</s>", "").replace("</s", "").replace("**", "").replace("*", "").strip()

    sentences = response.split('. ')
    unique = []
    seen = set()
    for s in sentences:
        norm = s.lower().strip()
        if norm not in seen and len(norm.split()) > 2:
            unique.append(s)
            seen.add(norm)
    return '. '.join(unique) if unique else response


def chat_reply(query, session_id=None):
    """
    Standalone chat function — no Flask dependency.
    Returns dict: {"response": str, "session_id": str}
    """
    if not query or not query.strip():
        return {"response": "Please provide a valid question or message.", "session_id": session_id or ""}

    if not session_id:
        session_id = str(uuid.uuid4())

    session = _get_session(session_id)

    # Redirect off-topic questions
    if not is_website_related(query):
        resp = "I'm here to help with questions about MAHE Innovation Centre. Please ask me about our events, resources, programs, or how to get involved with MiC."
        session["history"].append({"role": "user", "content": query})
        session["history"].append({"role": "assistant", "content": resp})
        return {"response": resp, "session_id": session_id}

    # Fallback when no Groq client
    if not client:
        resp = get_fallback_response(query)
        session["history"].append({"role": "user", "content": query})
        session["history"].append({"role": "assistant", "content": resp})
        return {"response": resp, "session_id": session_id}

    # Build messages for API
    messages = SystemChatBot.copy()
    messages.append({"role": "system", "content": get_realtime_information()})

    # Add recent history (last 10 messages)
    recent = session["history"][-10:]
    messages.extend(recent)
    messages.append({"role": "user", "content": query})

    try:
        completion = client.chat.completions.create(
            model="llama-3.3-70b-versatile",
            messages=messages,
            max_tokens=512,
            temperature=0.3,
            top_p=0.8,
            stream=True,
            stop=None
        )

        answer = ""
        for chunk in completion:
            if chunk.choices[0].delta.content:
                answer += chunk.choices[0].delta.content

        answer = clean_response(answer)

        if answer and len(answer) > 5:
            session["history"].append({"role": "user", "content": query})
            session["history"].append({"role": "assistant", "content": answer})
            # Trim history to last 20
            if len(session["history"]) > 20:
                session["history"] = session["history"][-20:]
            return {"response": answer, "session_id": session_id}
        else:
            return {"response": "I didn't generate a proper response. Please try rephrasing your question.", "session_id": session_id}

    except Exception as e:
        if "rate limit" in str(e).lower() or "429" in str(e):
            print("[CHATBOT] Rate limit hit, waiting 5s...")
            time.sleep(5)
            return chat_reply(query, session_id)
        print(f"[CHATBOT] Error: {e}")
        resp = get_fallback_response(query)
        return {"response": resp, "session_id": session_id}