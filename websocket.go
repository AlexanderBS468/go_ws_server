package main

import (
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/websocket"
)

const (
	websocketReadLimit  = 5120
	websocketPongWait   = 60 * time.Second
	websocketPingPeriod = 50 * time.Second
	websocketWriteWait  = 10 * time.Second
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func requireChannel(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "websocket channel is required", http.StatusBadRequest)
}

func (h *hub) handleWebSocket(w http.ResponseWriter, r *http.Request) {
	channel := strings.TrimPrefix(r.URL.Path, "/ws/")
	if channel == "" || channel == r.URL.Path {
		http.Error(w, "websocket channel is required", http.StatusBadRequest)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("websocket upgrade failed: %v", err)
		return
	}
	defer func() {
		h.remove(channel, conn)
		_ = conn.Close()
		h.publishUserCleanup(channel, conn.RemoteAddr().String())
	}()

	conn.SetReadLimit(websocketReadLimit)
	_ = conn.SetReadDeadline(time.Now().Add(websocketPongWait))
	conn.SetPongHandler(func(string) error {
		return conn.SetReadDeadline(time.Now().Add(websocketPongWait))
	})

	h.add(channel, conn)
	log.Printf("websocket connected: %s channel=%s", conn.RemoteAddr(), channel)

	done := make(chan struct{})
	defer close(done)
	go h.pingWebSocket(channel, conn, done)

	for {
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			log.Printf("websocket disconnected: %s: %v", conn.RemoteAddr(), err)
			return
		}
		if messageType != websocket.TextMessage {
			log.Printf("unsupported websocket message type from %s: %d", conn.RemoteAddr(), messageType)
			continue
		}

		payload, err := normalizeMessage(channel, conn.RemoteAddr().String(), message)
		if err != nil {
			log.Printf("invalid websocket message from %s: %v", conn.RemoteAddr(), err)
			continue
		}

		log.Printf("websocket message from %s channel=%s: %s", conn.RemoteAddr(), channel, payload)
		if err := h.broker.publish(channel, payload); err != nil {
			log.Printf("redis publish failed channel=%s: %v", channel, err)
		}
	}
}

func (h *hub) pingWebSocket(channel string, conn *websocket.Conn, done <-chan struct{}) {
	ticker := time.NewTicker(websocketPingPeriod)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if err := h.writePing(channel, conn); err != nil {
				log.Printf("websocket ping failed: %s: %v", conn.RemoteAddr(), err)
				_ = conn.Close()
				return
			}
		case <-done:
			return
		}
	}
}

func (h *hub) publishUserCleanup(channel string, from string) {
	payload, err := h.removeUser(from)
	if err != nil {
		log.Printf("user cleanup failed channel=%s: %v", channel, err)
		return
	}

	if err := h.broker.publish(channel, payload); err != nil {
		log.Printf("redis cleanup publish failed channel=%s: %v", channel, err)
	}
}
