package main

import (
	"encoding/json"
	"time"
)

type socketMessage struct {
	Event   string          `json:"event"`
	Channel string          `json:"channel"`
	Data    json.RawMessage `json:"data"`
	From    string          `json:"from,omitempty"`
	Ts      int64           `json:"ts"`
}

func normalizeMessage(channel string, from string, data []byte) ([]byte, error) {
	var msg socketMessage
	if err := json.Unmarshal(data, &msg); err != nil {
		return nil, err
	}

	if msg.Event == "" {
		msg.Event = "message"
	}
	if msg.Channel == "" {
		msg.Channel = channel
	}
	if msg.From == "" {
		msg.From = from
	}
	if msg.Ts == 0 {
		msg.Ts = time.Now().Unix()
	}

	return json.Marshal(msg)
}
