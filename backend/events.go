package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"

	"nhooyr.io/websocket"
)

// Event represents a presentation event sent to connected clients.
type Event struct {
	Type      string         `json:"type"`
	Source    string         `json:"source"`
	Timestamp int64          `json:"timestamp"`
	Data      map[string]any `json:"data,omitempty"`
}

// EventHub manages WebSocket clients and broadcasts events.
type EventHub struct {
	mu      sync.RWMutex
	clients map[*wsClient]struct{}
}

type wsClient struct {
	conn *websocket.Conn
	send chan []byte
}

func NewEventHub() *EventHub {
	return &EventHub{
		clients: make(map[*wsClient]struct{}),
	}
}

// Broadcast sends an event to all connected clients.
func (h *EventHub) Broadcast(evt Event) {
	if evt.Timestamp == 0 {
		evt.Timestamp = time.Now().UnixMilli()
	}
	data, err := json.Marshal(evt)
	if err != nil {
		log.Printf("event marshal error: %v", err)
		return
	}

	h.mu.RLock()
	defer h.mu.RUnlock()

	for c := range h.clients {
		select {
		case c.send <- data:
		default:
			// Client too slow, skip
		}
	}
}

// ServeHTTP handles WebSocket upgrade and client lifecycle.
func (h *EventHub) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	conn, err := websocket.Accept(w, r, &websocket.AcceptOptions{
		InsecureSkipVerify: true, // Allow connections from any origin in dev
	})
	if err != nil {
		log.Printf("ws accept error: %v", err)
		return
	}

	client := &wsClient{
		conn: conn,
		send: make(chan []byte, 64),
	}

	h.mu.Lock()
	h.clients[client] = struct{}{}
	clientCount := len(h.clients)
	h.mu.Unlock()

	log.Printf("event client connected (%d total)", clientCount)

	// Send a welcome event
	h.Broadcast(Event{
		Type:   "connection",
		Source: "hub",
		Data:   map[string]any{"clients": clientCount},
	})

	ctx := r.Context()

	// Writer goroutine: sends events to this client
	go func() {
		defer conn.CloseNow()
		for {
			select {
			case msg, ok := <-client.send:
				if !ok {
					return
				}
				err := conn.Write(ctx, websocket.MessageText, msg)
				if err != nil {
					return
				}
			case <-ctx.Done():
				return
			}
		}
	}()

	// Reader loop: handles incoming messages from the client (e.g., slide navigation events)
	for {
		_, msg, err := conn.Read(ctx)
		if err != nil {
			break
		}
		h.handleClientMessage(msg)
	}

	// Cleanup
	h.mu.Lock()
	delete(h.clients, client)
	clientCount = len(h.clients)
	h.mu.Unlock()
	close(client.send)

	log.Printf("event client disconnected (%d remaining)", clientCount)
}

// handleClientMessage processes messages sent from the frontend.
func (h *EventHub) handleClientMessage(msg []byte) {
	var evt Event
	if err := json.Unmarshal(msg, &evt); err != nil {
		log.Printf("invalid client event: %v", err)
		return
	}

	// Re-broadcast client events (e.g., slide-change) to all clients
	evt.Source = "client"
	if evt.Timestamp == 0 {
		evt.Timestamp = time.Now().UnixMilli()
	}
	h.Broadcast(evt)
}

// StartHeartbeat sends periodic heartbeat events so clients can detect disconnection.
func (h *EventHub) StartHeartbeat(ctx context.Context, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			h.Broadcast(Event{
				Type:   "heartbeat",
				Source: "hub",
			})
		case <-ctx.Done():
			return
		}
	}
}
