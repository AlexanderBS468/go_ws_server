package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("go-ws-server is running\n"))
	})
	mux.HandleFunc("/ws", handleWebSocket)

	addr := ":8080"
	log.Printf("server listening on %s", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatal(err)
	}
}

func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("websocket upgrade failed: %v", err)
		return
	}
	defer conn.Close()

	log.Printf("websocket connected: %s", conn.RemoteAddr())

	for {
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			log.Printf("websocket disconnected: %s: %v", conn.RemoteAddr(), err)
			return
		}

		log.Printf("websocket message from %s: %s", conn.RemoteAddr(), message)
		if err := conn.WriteMessage(messageType, message); err != nil {
			log.Printf("websocket write failed: %s: %v", conn.RemoteAddr(), err)
			return
		}
	}
}
