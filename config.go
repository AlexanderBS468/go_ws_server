package main

import (
	"os"
	"strings"
)

type config struct {
	ServerPort     string
	RedisAddr      string
	AllowedOrigins []string
}

func loadConfig() config {
	return config{
		ServerPort:     getEnv("SERVER_PORT", "8080"),
		RedisAddr:      getEnv("REDIS_ADDR", "localhost:6379"),
		AllowedOrigins: getEnvList("WS_ALLOWED_ORIGINS", []string{"http://localhost:8080", "http://127.0.0.1:8080"}),
	}
}

func getEnv(name string, fallback string) string {
	value := os.Getenv(name)
	if value == "" {
		return fallback
	}
	return value
}

func getEnvList(name string, fallback []string) []string {
	value := os.Getenv(name)
	if value == "" {
		return fallback
	}

	parts := strings.Split(value, ",")
	result := make([]string, 0, len(parts))
	for _, part := range parts {
		item := strings.TrimSpace(part)
		if item != "" {
			result = append(result, item)
		}
	}

	if len(result) == 0 {
		return fallback
	}
	return result
}
