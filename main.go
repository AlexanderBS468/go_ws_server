package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type hub struct {
	mu      sync.Mutex
	clients map[string]map[*websocket.Conn]struct{}
	broker  *redisBroker
}

type socketMessage struct {
	Event   string `json:"event"`
	Channel string `json:"channel"`
	Data    string `json:"data"`
	From    string `json:"from,omitempty"`
	Ts      int64  `json:"ts"`
}

type redisBroker struct {
	pool *redis.Pool
}

func newHub(broker *redisBroker) *hub {
	return &hub{
		clients: make(map[string]map[*websocket.Conn]struct{}),
		broker:  broker,
	}
}

func newRedisBroker(addr string) *redisBroker {
	return &redisBroker{
		pool: &redis.Pool{
			MaxIdle:   3,
			MaxActive: 10,
			Dial: func() (redis.Conn, error) {
				return redis.Dial("tcp", addr)
			},
		},
	}
}

func (b *redisBroker) publish(channel string, message []byte) error {
	conn := b.pool.Get()
	defer conn.Close()

	if err := conn.Err(); err != nil {
		return err
	}

	_, err := conn.Do("PUBLISH", channel, message)
	return err
}

func (b *redisBroker) subscribe(pattern string, handler func(channel string, message []byte)) error {
	conn := b.pool.Get()
	defer conn.Close()

	if err := conn.Err(); err != nil {
		return err
	}

	psc := redis.PubSubConn{Conn: conn}
	if err := psc.PSubscribe(pattern); err != nil {
		return err
	}

	for {
		switch event := psc.Receive().(type) {
		case redis.Message:
			handler(event.Channel, event.Data)
		case redis.Subscription:
			log.Printf("redis subscription: kind=%s channel=%s count=%d", event.Kind, event.Channel, event.Count)
		case error:
			return event
		}
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
	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		redisAddr = "localhost:6379"
	}

	broker := newRedisBroker(redisAddr)
	h := newHub(broker)
	go runRedisSubscriber(broker, h)

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

func runRedisSubscriber(broker *redisBroker, h *hub) {
	for {
		log.Printf("redis subscriber connecting")
		if err := broker.subscribe("*", h.broadcast); err != nil {
			log.Printf("redis subscriber stopped: %v", err)
			time.Sleep(2 * time.Second)
		}
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
		if err := h.broker.publish(channel, payload); err != nil {
			log.Printf("redis publish failed channel=%s: %v", channel, err)
		}
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
