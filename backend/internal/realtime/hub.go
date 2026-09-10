package realtime

import (
	"encoding/json"
	"sync"

	"github.com/gorilla/websocket"
)

// Hub keeps a registry of live websocket connections per user so services can push
// events (new notifications, new chat messages) without the client having to poll.
type Hub struct {
	mu    sync.RWMutex
	conns map[string]map[*websocket.Conn]bool
}

func NewHub() *Hub {
	return &Hub{conns: make(map[string]map[*websocket.Conn]bool)}
}
func (h *Hub) Register(userID string, conn *websocket.Conn) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if h.conns[userID] == nil {
		h.conns[userID] = make(map[*websocket.Conn]bool)
	}
	h.conns[userID][conn] = true
}
func (h *Hub) Unregister(userID string, conn *websocket.Conn) {
	h.mu.Lock()
	defer h.mu.Unlock()
	delete(h.conns[userID], conn)
	if len(h.conns[userID]) == 0 {
		delete(h.conns, userID)
	}
}

// Send delivers an event to every live connection for a user. Silently a no-op if the
// user has no open socket right now; callers should still persist state via the
// database as the source of truth (the socket is push-convenience, not the record).
func (h *Hub) Send(userID string, event any) {
	payload, err := json.Marshal(event)
	if err != nil {
		return
	}
	h.mu.RLock()
	conns := make([]*websocket.Conn, 0, len(h.conns[userID]))
	for conn := range h.conns[userID] {
		conns = append(conns, conn)
	}
	h.mu.RUnlock()
	for _, conn := range conns {
		if err := conn.WriteMessage(websocket.TextMessage, payload); err != nil {
			h.Unregister(userID, conn)
			_ = conn.Close()
		}
	}
}
