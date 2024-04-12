package main

import "os"

type config struct {
	ServerPort string
	RedisAddr  string
}

func loadConfig() config {
	return config{
		ServerPort: getEnv("SERVER_PORT", "8080"),
		RedisAddr:  getEnv("REDIS_ADDR", "localhost:6379"),
	}
}

func getEnv(name string, fallback string) string {
	value := os.Getenv(name)
	if value == "" {
		return fallback
	}
	return value
}
