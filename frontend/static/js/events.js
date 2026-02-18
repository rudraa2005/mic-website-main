async function fetchEvents() {
  try {
    const res = await fetch("/api/content/events/all");
    if (!res.ok) {
      console.error("Events fetch failed:", res.status);
      return [];
    }

    const data = await res.json();
    return Array.isArray(data) ? data : [];
  } catch (err) {
    console.error("Events fetch error:", err);
    return [];
  }
}
function formatDate(dateStr) {
  if (!dateStr) return "—";

  const d = new Date(dateStr);
  if (isNaN(d.getTime())) return "—";

  return d.toLocaleDateString("en-GB", {
    day: "2-digit",
    month: "short",
    year: "numeric"
  });
}
function renderEventCard(event) {
  const venue = event.venue || "TBA";
  const price = event.event_price || "TBA";
  console.log(price || "Nothing");
  console.log(event);

  return `
    <div class="event-card relative min-h-[360px] pb-6 max-w-md w-full mx-auto rounded-2xl overflow-hidden bg-gradient-to-br from-orange-primary to-orange-secondary shadow-xl">

      <div class="absolute inset-0 bg-black/30"></div>

      <div class="relative z-10 flex flex-col justify-between h-full p-6 text-white">

        <!-- Top -->
        <div>
          <p class="text-xs tracking-widest opacity-80">
            ${formatDate(event.event_date)}
          </p>

          <h3 class="mt-2 text-3xl font-black leading-tight">
            ${event.title}
          </h3>

          ${event.description
      ? `<p class="mt-3 text-sm text-white/90">${event.description}</p>`
      : ""
    }
        </div>

        <!-- Meta -->
        <div class="mt-6 space-y-3 text-sm text-white/90">
          <div class="flex items-center gap-3">
            <i class="fas fa-map-marker-alt w-4"></i>
            <span>${venue}</span>
          </div>

          <div class="flex items-center gap-3">
            <i class="fas fa-ticket-alt w-4"></i>
            <span>${event.price || 'Free'}</span>
          </div>

          <div class="flex items-center gap-3 text-white/70">
            <i class="fas fa-calendar-check w-4"></i>
            <span>${formatDate(event.event_date || "Upcoming")}</span>
          </div>
        </div>

        <!-- CTA -->
        <button class="self-start mt-6 px-5 py-2 rounded-full bg-white text-orange-primary font-semibold text-sm hover:bg-white/90 transition">
          <a
            href="${event.registration_link || '#'}"
            target="_blank"
            rel="noopener noreferrer"
            class="self-start mt-6  px-1 py-2 rounded-full bg-white text-orange-primary font-semibold text-sm hover:bg-white/90 transition"
          >
          Register Now →
        </a>
        </button>
      </div>
    </div>
  `;
}

async function loadEvents() {
  const container = document.getElementById("eventsGrid");
  if (!container) return;

  const events = await fetchEvents();

  if (!Array.isArray(events) || events.length === 0) {
    container.innerHTML =
      `<p class="text-gray-600">No events available right now.</p>`;
    return;
  }

  container.innerHTML = events
    .sort((a, b) => {
      const da = new Date(a.content_data?.date || 0);
      const db = new Date(b.content_data?.date || 0);
      return da - db;
    })
    .map(renderEventCard)
    .join("")
}

document.addEventListener("DOMContentLoaded", loadEvents);