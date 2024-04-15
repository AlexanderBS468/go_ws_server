package main

import (
	"encoding/json"
	"sort"
)

type userState struct {
	Hash     string
	ID       int
	Fullname string
	Modals   []int
}

type loginPayload struct {
	Hash     string `json:"hash"`
	ID       int    `json:"id"`
	Fullname string `json:"fullname"`
}

type leadPayload struct {
	LeadID int `json:"leadId"`
}

func (h *hub) processMessage(message []byte) ([]byte, error) {
	var msg socketMessage
	if err := json.Unmarshal(message, &msg); err != nil {
		return nil, err
	}

	switch msg.Event {
	case "login":
		return h.processLogin(msg)
	case "lead-details-open":
		return h.processLeadOpen(msg)
	case "lead-details-close":
		return h.processLeadClose(msg)
	case "lead-page-out":
		return h.processLeadPageOut(msg)
	default:
		return nil, nil
	}
}

func (h *hub) processLogin(msg socketMessage) ([]byte, error) {
	var payload loginPayload
	if err := json.Unmarshal(msg.Data, &payload); err != nil {
		return nil, err
	}

	h.mu.Lock()
	h.users[msg.From] = userState{
		Hash:     payload.Hash,
		ID:       payload.ID,
		Fullname: payload.Fullname,
		Modals:   []int{},
	}
	h.mu.Unlock()

	return h.modalsInfoMessage()
}

func (h *hub) processLeadOpen(msg socketMessage) ([]byte, error) {
	var payload leadPayload
	if err := json.Unmarshal(msg.Data, &payload); err != nil {
		return nil, err
	}

	h.mu.Lock()
	user := h.users[msg.From]
	if !containsInt(user.Modals, payload.LeadID) {
		user.Modals = append(user.Modals, payload.LeadID)
	}
	h.users[msg.From] = user
	h.mu.Unlock()

	return h.modalsInfoMessage()
}

func (h *hub) processLeadClose(msg socketMessage) ([]byte, error) {
	var payload leadPayload
	if err := json.Unmarshal(msg.Data, &payload); err != nil {
		return nil, err
	}

	h.mu.Lock()
	user := h.users[msg.From]
	user.Modals = removeInt(user.Modals, payload.LeadID)
	h.users[msg.From] = user
	h.mu.Unlock()

	return h.modalsInfoMessage()
}

func (h *hub) processLeadPageOut(msg socketMessage) ([]byte, error) {
	h.mu.Lock()
	user := h.users[msg.From]
	user.Modals = []int{}
	h.users[msg.From] = user
	h.mu.Unlock()

	return h.modalsInfoMessage()
}

func (h *hub) modalsInfoMessage() ([]byte, error) {
	data, err := json.Marshal(h.userModals())
	if err != nil {
		return nil, err
	}

	msg := socketMessage{
		Event: "modalsInfo",
		Data:  data,
	}

	return json.Marshal(msg)
}

func (h *hub) userModals() map[string][]int {
	h.mu.Lock()
	defer h.mu.Unlock()

	result := make(map[string][]int)
	for _, user := range h.users {
		if user.Fullname == "" {
			continue
		}

		for _, leadID := range user.Modals {
			if !containsInt(result[user.Fullname], leadID) {
				result[user.Fullname] = append(result[user.Fullname], leadID)
			}
		}
	}

	for name := range result {
		sort.Ints(result[name])
	}

	return result
}

func containsInt(items []int, value int) bool {
	for _, item := range items {
		if item == value {
			return true
		}
	}
	return false
}

func removeInt(items []int, value int) []int {
	filtered := make([]int, 0, len(items))
	for _, item := range items {
		if item != value {
			filtered = append(filtered, item)
		}
	}
	return filtered
}
