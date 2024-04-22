package main

import "log"

func main() {
	cfg := loadConfig()

	ctx, stop := shutdownContext()
	defer stop()

	broker := newRedisBroker(cfg.RedisAddr)
	h := newHub(broker)
	go runRedisSubscriber(ctx, broker, h)

	if err := runHTTPServer(ctx, cfg, newRouter(h)); err != nil {
		log.Fatal(err)
	}
}
