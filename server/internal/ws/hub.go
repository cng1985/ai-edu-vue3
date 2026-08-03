package ws

import (
	"encoding/json"
	"sync"

	"github.com/gorilla/websocket"
	"go.uber.org/fx"
)

type Client struct {
	UserID   string
	Username string
	Role     string
	IsAdmin  bool
	Conn     *websocket.Conn
	Send     chan []byte
	Tickets  map[string]bool
}

type Hub struct {
	mu      sync.RWMutex
	clients map[*Client]bool
}

type Envelope struct {
	Event   string          `json:"event"`
	Payload json.RawMessage `json:"payload"`
}

func NewHub() *Hub {
	return &Hub{clients: make(map[*Client]bool)}
}

var Module = fx.Provide(NewHub)

func (h *Hub) Register(c *Client) {
	h.mu.Lock()
	h.clients[c] = true
	h.mu.Unlock()
}

func (h *Hub) Unregister(c *Client) {
	h.mu.Lock()
	if _, ok := h.clients[c]; ok {
		delete(h.clients, c)
		close(c.Send)
	}
	h.mu.Unlock()
}

func (h *Hub) Subscribe(c *Client, ticketID string) {
	h.mu.Lock()
	if c.Tickets == nil {
		c.Tickets = make(map[string]bool)
	}
	c.Tickets[ticketID] = true
	h.mu.Unlock()
}

func (h *Hub) BroadcastToTicket(ticketID string, data []byte, exclude *Client) {
	h.mu.RLock()
	defer h.mu.RUnlock()
	for c := range h.clients {
		if exclude != nil && c == exclude {
			continue
		}
		if c.Tickets[ticketID] || c.IsAdmin {
			select {
			case c.Send <- data:
			default:
			}
		}
	}
}

func (h *Hub) BroadcastToAdmins(data []byte) {
	h.mu.RLock()
	defer h.mu.RUnlock()
	for c := range h.clients {
		if c.IsAdmin {
			select {
			case c.Send <- data:
			default:
			}
		}
	}
}

func (h *Hub) SendToUser(userID string, data []byte) {
	h.mu.RLock()
	defer h.mu.RUnlock()
	for c := range h.clients {
		if c.UserID == userID {
			select {
			case c.Send <- data:
			default:
			}
		}
	}
}

func MarshalEvent(event string, payload any) ([]byte, error) {
	p, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	return json.Marshal(Envelope{Event: event, Payload: p})
}
