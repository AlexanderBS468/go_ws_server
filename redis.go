package main

import (
	"context"
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

func (b *redisBroker) ping() error {
	conn := b.pool.Get()
	defer conn.Close()

	if err := conn.Err(); err != nil {
		return err
	}

	_, err := conn.Do("PING")
	return err
}

func (b *redisBroker) subscribe(ctx context.Context, pattern string, handler func(channel string, message []byte)) error {
	conn := b.pool.Get()
	defer conn.Close()

	if err := conn.Err(); err != nil {
		return err
	}

	psc := redis.PubSubConn{Conn: conn}
	if err := psc.PSubscribe(pattern); err != nil {
		return err
	}

	done := make(chan struct{})
	defer close(done)

	go func() {
		select {
		case <-ctx.Done():
			_ = conn.Close()
		case <-done:
		}
	}()

	for {
		switch event := psc.Receive().(type) {
		case redis.Message:
			handler(event.Channel, event.Data)
		case redis.Subscription:
			log.Printf("redis subscription: kind=%s channel=%s count=%d", event.Kind, event.Channel, event.Count)
		case error:
			if ctx.Err() != nil {
				return ctx.Err()
			}
			return event
		}
	}
}

func runRedisSubscriber(ctx context.Context, broker *redisBroker, h *hub) {
	for {
		select {
		case <-ctx.Done():
			log.Printf("redis subscriber stopped")
			return
		default:
		}

		log.Printf("redis subscriber connecting")
		if err := broker.subscribe(ctx, "*", h.handleRedisMessage); err != nil {
			if ctx.Err() != nil {
				log.Printf("redis subscriber stopped")
				return
			}
			log.Printf("redis subscriber stopped: %v", err)
			select {
			case <-time.After(2 * time.Second):
			case <-ctx.Done():
				log.Printf("redis subscriber stopped")
				return
			}
		}
	}
}
