package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func shutdownContext() (context.Context, context.CancelFunc) {
	return signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
}

func newRouter(h *hub) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("go-ws-server is running\n"))
	})
	mux.HandleFunc("/ws", requireChannel)
	mux.HandleFunc("/ws/", h.handleWebSocket)
	mux.Handle("/", http.FileServer(http.Dir("./public")))

	return mux
}

func runHTTPServer(ctx context.Context, cfg config, handler http.Handler) error {
	addr := ":" + cfg.ServerPort
	server := &http.Server{
		Addr:              addr,
		Handler:           handler,
		ReadHeaderTimeout: 5 * time.Second,
	}

	serverErrors := make(chan error, 1)
	go func() {
		serverErrors <- server.ListenAndServe()
	}()

	log.Printf("server listening on %s redis=%s", addr, cfg.RedisAddr)

	select {
	case err := <-serverErrors:
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			return err
		}
	case <-ctx.Done():
		log.Printf("shutdown signal received")

		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer shutdownCancel()

		if err := server.Shutdown(shutdownCtx); err != nil {
			log.Printf("server graceful shutdown failed: %v", err)
			if closeErr := server.Close(); closeErr != nil {
				log.Printf("server close failed: %v", closeErr)
			}
		}
	}

	return nil
}
