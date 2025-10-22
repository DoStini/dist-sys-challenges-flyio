package main

import (
	"context"
	"encoding/json"
	"fmt"
	"sync/atomic"

	"github.com/dostini/dist-sys-challenges-flyio/pkg/node"
)

type broadcastRequest struct {
	Message int `json:"message"`
}

type gossipRequest struct {
	OriginNode string `json:"origin_node"`
	SeqID      int    `json:"seq_id"`
	Message    int    `json:"message"`
	Type       string `json:"type"`
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

	sequenceSet := NewAtomicMap[bool]()
	sequenceCounter := atomic.Int64{}

	handleMessage := func(msg node.Message, req gossipRequest) error {
		key := fmt.Sprintf("%s-%d", req.OriginNode, req.SeqID)
		if _, found := sequenceSet.GetOrDefault(key, true); found {
			// Ignore already processed message
			return nil
		}

		data = append(data, req.Message)
		for _, nodeId := range n.Neighbours() {
			if nodeId == msg.Src || nodeId == req.OriginNode {
				continue
			}

			n.Send(nodeId, req)
		}

		return nil
	}

	n.Handle("broadcast", func(msg node.Message) error {
		var req broadcastRequest
		json.Unmarshal(msg.Body, &req)

		gossipRequest := gossipRequest{
			OriginNode: n.ID(),
			SeqID:      int(sequenceCounter.Add(1)),
			Message:    req.Message,
			Type:       "gossip",
		}

		handleMessage(msg, gossipRequest)

		return n.Reply(msg, map[string]string{"type": "broadcast_ok"})
	})

	n.Handle("gossip", func(msg node.Message) error {
		var req gossipRequest
		json.Unmarshal(msg.Body, &req)

		err := handleMessage(msg, req)
		if err != nil {
			return err
		}

		return nil
	})

	n.Handle("read", func(msg node.Message) error {
		return n.Reply(msg, newReadResponse(data))
	})

	n.Handle("topology", func(msg node.Message) error {
		return n.Reply(msg, map[string]string{"type": "topology_ok"})
	})

	n.Run()
}
