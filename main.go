package main

import (
	"log"
	"net/http"
)

func main() {
	cfg := loadConfig()

	broker := newRedisBroker(cfg.RedisAddr)
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

	addr := ":" + cfg.ServerPort
	log.Printf("server listening on %s redis=%s", addr, cfg.RedisAddr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatal(err)
	}
}
