package main

import (
	"net/http/httptest"
	"testing"
)

func TestCheckWebSocketOriginAllowsConfiguredOrigin(t *testing.T) {
	h := newHub(nil, []string{"http://localhost:8080"})
	r := httptest.NewRequest("GET", "/ws/leads", nil)
	r.Header.Set("Origin", "http://localhost:8080")

	if !h.checkWebSocketOrigin(r) {
		t.Fatal("expected configured origin to be allowed")
	}
}

func TestCheckWebSocketOriginRejectsUnknownOrigin(t *testing.T) {
	h := newHub(nil, []string{"http://localhost:8080"})
	r := httptest.NewRequest("GET", "/ws/leads", nil)
	r.Header.Set("Origin", "https://example.com")

	if h.checkWebSocketOrigin(r) {
		t.Fatal("expected unknown origin to be rejected")
	}
}

func TestCheckWebSocketOriginAllowsMissingOrigin(t *testing.T) {
	h := newHub(nil, []string{"http://localhost:8080"})
	r := httptest.NewRequest("GET", "/ws/leads", nil)

	if !h.checkWebSocketOrigin(r) {
		t.Fatal("expected missing origin to be allowed")
	}
}
