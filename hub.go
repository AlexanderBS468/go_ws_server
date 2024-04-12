package main

import (
	"log"
	"sync"

	"github.com/gorilla/websocket"
)

type hub struct {
	mu      sync.Mutex
	clients map[string]map[*websocket.Conn]struct{}
	broker  *redisBroker
}

func newHub(broker *redisBroker) *hub {
	return &hub{
		clients: make(map[string]map[*websocket.Conn]struct{}),
		broker:  broker,
	}
}

func (h *hub) add(channel string, conn *websocket.Conn) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if h.clients[channel] == nil {
		h.clients[channel] = make(map[*websocket.Conn]struct{})
	}
	h.clients[channel][conn] = struct{}{}
}

func (h *hub) remove(channel string, conn *websocket.Conn) {
	h.mu.Lock()
	defer h.mu.Unlock()

	delete(h.clients[channel], conn)
	if len(h.clients[channel]) == 0 {
		delete(h.clients, channel)
	}
}

func (h *hub) broadcast(channel string, message []byte) {
	h.mu.Lock()
	defer h.mu.Unlock()

	for conn := range h.clients[channel] {
		if err := conn.WriteMessage(websocket.TextMessage, message); err != nil {
			log.Printf("websocket broadcast failed: %s: %v", conn.RemoteAddr(), err)
			_ = conn.Close()
			delete(h.clients[channel], conn)
		}
	}
}
