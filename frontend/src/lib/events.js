import { writable, derived } from 'svelte/store';

// All events received from the server
export const events = writable([]);

// The most recent event
export const latestEvent = derived(events, ($events) =>
  $events.length > 0 ? $events[$events.length - 1] : null
);

// Filter events by type
export function eventsByType(type) {
  return derived(events, ($events) =>
    $events.filter((e) => e.type === type)
  );
}

// The most recent event of a specific type
export function latestEventOfType(type) {
  return derived(events, ($events) => {
    for (let i = $events.length - 1; i >= 0; i--) {
      if ($events[i].type === type) return $events[i];
    }
    return null;
  });
}

// Commands derived from events
export const commands = derived(events, ($events) =>
  $events
    .filter((e) => e.type === 'command' && e.data?.command)
    .map((e) => ({
      command: e.data.command,
      exitCode: e.data.exit_code,
      timestamp: e.timestamp,
    }))
);

// The latest command
export const latestCommand = derived(commands, ($commands) =>
  $commands.length > 0 ? $commands[$commands.length - 1] : null
);

// Connection state
export const connected = writable(false);

// Max events to keep in memory
const MAX_EVENTS = 500;

let ws = null;
let reconnectTimer = null;

export function connect() {
  if (ws) return;

  const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
  const url = `${protocol}//${window.location.host}/ws/events`;

  ws = new WebSocket(url);

  ws.onopen = () => {
    connected.set(true);
    console.log('[events] connected');
  };

  ws.onmessage = (msg) => {
    try {
      const evt = JSON.parse(msg.data);
      // Skip heartbeats from the event list
      if (evt.type === 'heartbeat') return;

      events.update(($events) => {
        const next = [...$events, evt];
        // Trim if over limit
        if (next.length > MAX_EVENTS) {
          return next.slice(next.length - MAX_EVENTS);
        }
        return next;
      });
    } catch (e) {
      console.warn('[events] invalid message:', msg.data);
    }
  };

  ws.onclose = () => {
    connected.set(false);
    ws = null;
    console.log('[events] disconnected, reconnecting in 2s');
    reconnectTimer = setTimeout(connect, 2000);
  };

  ws.onerror = (err) => {
    console.warn('[events] error:', err);
    ws?.close();
  };
}

export function disconnect() {
  clearTimeout(reconnectTimer);
  ws?.close();
  ws = null;
  connected.set(false);
}

// Send an event to the server (e.g., slide navigation)
export function send(type, data = {}) {
  if (!ws || ws.readyState !== WebSocket.OPEN) return;
  ws.send(JSON.stringify({ type, data }));
}

// Clear event history
export function clearEvents() {
  events.set([]);
}
