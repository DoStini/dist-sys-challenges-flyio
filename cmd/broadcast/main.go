package main

import (
	"encoding/json"
	"log"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

type broadcastRequest struct {
	Message int `json:"message"`
}

type readResponse struct {
	Type     string `json:"type"`
	Messages []int  `json:"messages"`
}

func newReadResponse(data []int) readResponse {
	return readResponse{
		Type:     "read_ok",
		Messages: data,
	}
}

func main() {
	n := maelstrom.NewNode()

	var data []int

	n.Handle("broadcast", func(msg maelstrom.Message) error {
		var req broadcastRequest
		json.Unmarshal(msg.Body, &req)

		data = append(data, req.Message)

		return n.Reply(msg, map[string]string{"type": "broadcast_ok"})
	})

	n.Handle("read", func(msg maelstrom.Message) error {
		return n.Reply(msg, newReadResponse(data))
	})

	n.Handle("topology", func(msg maelstrom.Message) error {
		return n.Reply(msg, map[string]string{"type": "topology_ok"})
	})

	if err := n.Run(); err != nil {
		log.Fatal(err)
	}
}
