package main

import (
	"encoding/json"
	"testing"
)

func TestLeadModalStateEvents(t *testing.T) {
	h := newHub(nil)

	processTestMessage(t, h, socketMessage{
		Event: "login",
		From:  "browser-1",
		Data:  mustRawJSON(t, loginPayload{Hash: "browser-1", ID: 7, Fullname: "Alice"}),
	})
	processTestMessage(t, h, socketMessage{
		Event: "lead-details-open",
		From:  "browser-1",
		Data:  mustRawJSON(t, leadPayload{LeadID: 101}),
	})
	processTestMessage(t, h, socketMessage{
		Event: "lead-details-open",
		From:  "browser-1",
		Data:  mustRawJSON(t, leadPayload{LeadID: 102}),
	})
	processTestMessage(t, h, socketMessage{
		Event: "lead-details-open",
		From:  "browser-1",
		Data:  mustRawJSON(t, leadPayload{LeadID: 101}),
	})

	assertModals(t, h.userModals(), map[string][]int{
		"Alice": {101, 102},
	})

	processTestMessage(t, h, socketMessage{
		Event: "lead-details-close",
		From:  "browser-1",
		Data:  mustRawJSON(t, leadPayload{LeadID: 101}),
	})

	assertModals(t, h.userModals(), map[string][]int{
		"Alice": {102},
	})

	processTestMessage(t, h, socketMessage{
		Event: "lead-page-out",
		From:  "browser-1",
		Data:  mustRawJSON(t, map[string]string{}),
	})

	assertModals(t, h.userModals(), map[string][]int{})
}

func TestRemoveUserClearsModalState(t *testing.T) {
	h := newHub(nil)

	processTestMessage(t, h, socketMessage{
		Event: "login",
		From:  "browser-1",
		Data:  mustRawJSON(t, loginPayload{Hash: "browser-1", ID: 7, Fullname: "Alice"}),
	})
	processTestMessage(t, h, socketMessage{
		Event: "lead-details-open",
		From:  "browser-1",
		Data:  mustRawJSON(t, leadPayload{LeadID: 101}),
	})

	if _, err := h.removeUser("browser-1"); err != nil {
		t.Fatalf("remove user: %v", err)
	}

	assertModals(t, h.userModals(), map[string][]int{})
}

func processTestMessage(t *testing.T, h *hub, msg socketMessage) {
	t.Helper()

	data, err := json.Marshal(msg)
	if err != nil {
		t.Fatalf("marshal message: %v", err)
	}

	if _, err := h.processMessage(data); err != nil {
		t.Fatalf("process message %s: %v", msg.Event, err)
	}
}

func mustRawJSON(t *testing.T, value interface{}) json.RawMessage {
	t.Helper()

	data, err := json.Marshal(value)
	if err != nil {
		t.Fatalf("marshal raw json: %v", err)
	}

	return data
}

func assertModals(t *testing.T, got map[string][]int, want map[string][]int) {
	t.Helper()

	gotData, err := json.Marshal(got)
	if err != nil {
		t.Fatalf("marshal got: %v", err)
	}
	wantData, err := json.Marshal(want)
	if err != nil {
		t.Fatalf("marshal want: %v", err)
	}

	if string(gotData) != string(wantData) {
		t.Fatalf("modals mismatch: got %s want %s", gotData, wantData)
	}
}
