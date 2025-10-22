package main

import (
	"context"
	"encoding/json"

	"github.com/dostini/dist-sys-challenges-flyio/pkg/node"
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
	ctx := context.Background()
	n := node.NewNode(ctx)

	var data []int

	n.Handle("broadcast", func(msg node.Message) error {
		var req broadcastRequest
		json.Unmarshal(msg.Body, &req)

		data = append(data, req.Message)

		return n.Reply(msg, map[string]string{"type": "broadcast_ok"})
	})

	n.Handle("read", func(msg node.Message) error {
		return n.Reply(msg, newReadResponse(data))
	})

	n.Handle("topology", func(msg node.Message) error {
		return n.Reply(msg, map[string]string{"type": "topology_ok"})
	})

	n.Run()
}
