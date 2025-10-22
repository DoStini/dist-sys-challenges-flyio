package main

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"sync/atomic"

	"github.com/dostini/dist-sys-challenges-flyio/pkg/node"
)

type GenericMessage struct {
	Type string `json:"type"`
}

type MaxResponse struct {
	GenericMessage
	Max int `json:"max"`
}

func NewMaxResponse(max int) MaxResponse {
	return MaxResponse{
		GenericMessage: GenericMessage{
			Type: "max_ok",
		},
		Max: max,
	}
}

func main() {
	ctx := context.Background()
	n := node.NewNode(ctx)

	currMaxFirstSec := atomic.Uint64{}
	currMaxSecondSec := atomic.Uint64{}

	n.Handle("generate", func(msg node.Message) error {
		// Unmarshal the message body as an loosely-typed map.
		var body map[string]any
		if err := json.Unmarshal(msg.Body, &body); err != nil {
			return err
		}

		maxNumFirstSec := currMaxFirstSec.Add(1)
		maxNumSecondSec := currMaxSecondSec.Load()
		if maxNumFirstSec == math.MaxUint64 {
			currMaxFirstSec.Store(0)
			maxNumSecondSec = currMaxSecondSec.Add(1)
		}

		// Update the message type to return back.
		body["type"] = "generate_ok"
		body["id"] = fmt.Sprintf("%s-%d-%d", n.ID(), maxNumFirstSec, maxNumSecondSec)

		return n.Reply(msg, body)
	})

	n.Run()
}
