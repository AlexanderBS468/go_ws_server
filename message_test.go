package main

import (
	"encoding/json"
	"testing"
	"time"
)

func TestNormalizeMessageFillsMissingFields(t *testing.T) {
	before := time.Now().Unix()

	out, err := normalizeMessage("leads", "browser-1", []byte(`{"data":{"text":"hello"}}`))
	if err != nil {
		t.Fatalf("normalize message: %v", err)
	}

	after := time.Now().Unix()
	msg := decodeSocketMessage(t, out)

	if msg.Event != "message" {
		t.Fatalf("event mismatch: got %q", msg.Event)
	}
	if msg.Channel != "leads" {
		t.Fatalf("channel mismatch: got %q", msg.Channel)
	}
	if msg.From != "browser-1" {
		t.Fatalf("from mismatch: got %q", msg.From)
	}
	if msg.Ts < before || msg.Ts > after {
		t.Fatalf("timestamp out of range: got %d want between %d and %d", msg.Ts, before, after)
	}
	if string(msg.Data) != `{"text":"hello"}` {
		t.Fatalf("data mismatch: got %s", msg.Data)
	}
}

func TestNormalizeMessageKeepsExistingFields(t *testing.T) {
	input := []byte(`{"event":"login","channel":"users","from":"browser-2","ts":1689588123,"data":{"id":7}}`)

	out, err := normalizeMessage("leads", "browser-1", input)
	if err != nil {
		t.Fatalf("normalize message: %v", err)
	}

	msg := decodeSocketMessage(t, out)

	if msg.Event != "login" {
		t.Fatalf("event mismatch: got %q", msg.Event)
	}
	if msg.Channel != "users" {
		t.Fatalf("channel mismatch: got %q", msg.Channel)
	}
	if msg.From != "browser-2" {
		t.Fatalf("from mismatch: got %q", msg.From)
	}
	if msg.Ts != 1689588123 {
		t.Fatalf("timestamp mismatch: got %d", msg.Ts)
	}
	if string(msg.Data) != `{"id":7}` {
		t.Fatalf("data mismatch: got %s", msg.Data)
	}
}

func TestNormalizeMessageReturnsInvalidJSONError(t *testing.T) {
	if _, err := normalizeMessage("leads", "browser-1", []byte(`{bad json}`)); err == nil {
		t.Fatal("expected invalid JSON error")
	}
}

func decodeSocketMessage(t *testing.T, data []byte) socketMessage {
	t.Helper()

	var msg socketMessage
	if err := json.Unmarshal(data, &msg); err != nil {
		t.Fatalf("decode socket message: %v", err)
	}

	return msg
}
