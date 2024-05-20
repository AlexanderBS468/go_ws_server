package main

import (
	"log"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type hub struct {
	mu             sync.Mutex
	clients        map[string]map[*websocket.Conn]struct{}
	users          map[string]userState
	broker         *redisBroker
	allowedOrigins map[string]struct{}
}

func newHub(broker *redisBroker, allowedOrigins []string) *hub {
	return &hub{
		clients:        make(map[string]map[*websocket.Conn]struct{}),
		users:          make(map[string]userState),
		broker:         broker,
		allowedOrigins: originSet(allowedOrigins),
	}
}

func originSet(origins []string) map[string]struct{} {
	result := make(map[string]struct{}, len(origins))
	for _, origin := range origins {
		result[origin] = struct{}{}
	}
	return result
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

func (h *hub) writePing(channel string, conn *websocket.Conn) error {
	h.mu.Lock()
	defer h.mu.Unlock()

	if _, ok := h.clients[channel][conn]; !ok {
		return nil
	}

	return conn.WriteControl(websocket.PingMessage, nil, time.Now().Add(websocketWriteWait))
}

func (h *hub) handleRedisMessage(channel string, message []byte) {
	out, err := h.processMessage(message)
	if err != nil {
		log.Printf("event processing failed channel=%s: %v", channel, err)
		return
	}
	if out == nil {
		out = message
	}

	h.broadcast(channel, out)
}
