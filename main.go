package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type hub struct {
	mu      sync.Mutex
	clients map[string]map[*websocket.Conn]struct{}
}

type socketMessage struct {
	Event   string `json:"event"`
	Channel string `json:"channel"`
	Data    string `json:"data"`
	From    string `json:"from,omitempty"`
	Ts      int64  `json:"ts"`
}

func newHub() *hub {
	return &hub{
		clients: make(map[string]map[*websocket.Conn]struct{}),
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

func main() {
	h := newHub()
	mux := http.NewServeMux()
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("go-ws-server is running\n"))
	})
	mux.HandleFunc("/ws", requireChannel)
	mux.HandleFunc("/ws/", h.handleWebSocket)
	mux.Handle("/", http.FileServer(http.Dir("./public")))

	addr := ":8080"
	log.Printf("server listening on %s", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatal(err)
	}
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
	}()

	h.add(channel, conn)
	log.Printf("websocket connected: %s channel=%s", conn.RemoteAddr(), channel)

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
		h.broadcast(channel, payload)
	}
}

func normalizeMessage(channel string, from string, data []byte) ([]byte, error) {
	var msg socketMessage
	if err := json.Unmarshal(data, &msg); err != nil {
		return nil, err
	}

	if msg.Event == "" {
		msg.Event = "message"
	}
	if msg.Channel == "" {
		msg.Channel = channel
	}
	if msg.From == "" {
		msg.From = from
	}
	if msg.Ts == 0 {
		msg.Ts = time.Now().Unix()
	}

	return json.Marshal(msg)
}
