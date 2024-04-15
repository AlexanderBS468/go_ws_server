package main

import (
	"log"
	"time"

	"github.com/gomodule/redigo/redis"
)

type redisBroker struct {
	pool *redis.Pool
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

func runRedisSubscriber(broker *redisBroker, h *hub) {
	for {
		log.Printf("redis subscriber connecting")
		if err := broker.subscribe("*", h.handleRedisMessage); err != nil {
			log.Printf("redis subscriber stopped: %v", err)
			time.Sleep(2 * time.Second)
		}
	}
}
