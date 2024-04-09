package main

import (
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type hub struct {
	mu      sync.Mutex
	clients map[*websocket.Conn]struct{}
}

func newHub() *hub {
	return &hub{
		clients: make(map[*websocket.Conn]struct{}),
	}
}

func (h *hub) add(conn *websocket.Conn) {
	h.mu.Lock()
	defer h.mu.Unlock()

	h.clients[conn] = struct{}{}
}

func (h *hub) remove(conn *websocket.Conn) {
	h.mu.Lock()
	defer h.mu.Unlock()

	delete(h.clients, conn)
}

func (h *hub) broadcast(messageType int, message []byte) {
	h.mu.Lock()
	defer h.mu.Unlock()

	for conn := range h.clients {
		if err := conn.WriteMessage(messageType, message); err != nil {
			log.Printf("websocket broadcast failed: %s: %v", conn.RemoteAddr(), err)
			_ = conn.Close()
			delete(h.clients, conn)
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
	mux.HandleFunc("/ws", h.handleWebSocket)
	mux.Handle("/", http.FileServer(http.Dir("./public")))

	addr := ":8080"
	log.Printf("server listening on %s", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatal(err)
	}
}

func (h *hub) handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("websocket upgrade failed: %v", err)
		return
	}
	defer func() {
		h.remove(conn)
		_ = conn.Close()
	}()

	h.add(conn)
	log.Printf("websocket connected: %s", conn.RemoteAddr())

	for {
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			log.Printf("websocket disconnected: %s: %v", conn.RemoteAddr(), err)
			return
		}

		log.Printf("websocket message from %s: %s", conn.RemoteAddr(), message)
		h.broadcast(messageType, message)
	}
}
