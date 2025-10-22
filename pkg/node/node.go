package node

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"sync/atomic"
)

type Message struct {
	Src  string          `json:"src,omitempty"`
	Dest string          `json:"dest,omitempty"`
	Body json.RawMessage `json:"body,omitempty"`
}

type Node struct {
	ctx context.Context

	id         string
	neighbours []string

	inboundChan  chan Message
	outboundChan chan Message

	closeInboundChan  chan struct{}
	closeOutboundChan chan struct{}

	idCounter atomic.Int64

	handlers map[string]func(Message) error
}

func NewNode(ctx context.Context) *Node {
	n := Node{
		ctx:               ctx,
		inboundChan:       make(chan Message, 64),
		outboundChan:      make(chan Message, 64),
		closeInboundChan:  make(chan struct{}),
		closeOutboundChan: make(chan struct{}),
		handlers:          map[string]func(Message) error{},
		idCounter:         atomic.Int64{},
	}

	n.Handle("init", n.init)

	return &n
}

func (n *Node) init(msg Message) error {
	request := InitRequest{}
	err := json.Unmarshal(msg.Body, &request)
	if err != nil {
		return err
	}

	n.id = request.NodeID
	n.neighbours = request.NodeIDs

	return n.Reply(msg, map[string]string{
		"type": "init_ok",
	})
}

func (n *Node) Run() {
	go n.run()
	go n.stdout()
	go n.stdin()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	// Block until a signal is received.
	<-c

	slog.Info("shutting down node")
	n.Close()

	slog.Info("shutted down node")
}

func (n *Node) Close() {
	n.closeInboundChan <- struct{}{}
	n.closeOutboundChan <- struct{}{}

	<-n.closeInboundChan
	<-n.closeOutboundChan

	close(n.closeInboundChan)
	close(n.closeOutboundChan)

	slog.Info("node closed sucessfully")
}

func (n *Node) Handle(messageType string, callback func(Message) error) {
	n.handlers[messageType] = callback
}

func (n *Node) Send(dest string, body any) error {
	msgId := n.idCounter.Add(1)

	tempJson, err := json.Marshal(body)
	if err != nil {
		return err
	}

	var data map[string]any
	err = json.Unmarshal(tempJson, &data)
	if err != nil {
		return err
	}

	data["msg_id"] = msgId

	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	msg := Message{
		Src:  n.id,
		Dest: dest,
		Body: jsonData,
	}

	n.outboundChan <- msg
	slog.Info("sent message", "message", msg)

	return nil
}

func (n *Node) Reply(msg Message, body any) error {
	var parsedPayload GenericPayload

	err := json.Unmarshal(msg.Body, &parsedPayload)
	if err != nil {
		return err
	}

	tempJson, err := json.Marshal(body)
	if err != nil {
		return err
	}

	var data map[string]any
	err = json.Unmarshal(tempJson, &data)
	if err != nil {
		return err
	}

	data["in_reply_to"] = parsedPayload.MsgID

	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	src := msg.Src
	msg.Src = msg.Dest
	msg.Dest = src
	msg.Body = jsonData

	n.outboundChan <- msg
	slog.Info("replied message", "message", msg)

	return nil
}

func (n *Node) ID() string {
	return n.id
}

func (n *Node) Neighbours() []string {
	return n.neighbours
}

func (n *Node) handleMessage(message Message) error {
	type genericMsg struct {
		Type string `json:"type"`
	}

	var parsed genericMsg

	err := json.Unmarshal(message.Body, &parsed)
	if err != nil {
		return err
	}

	return n.handlers[parsed.Type](message)
}

func (n *Node) stdin() {
	defer func() {
		n.closeInboundChan <- struct{}{}
		close(n.inboundChan)
		slog.Info("closed succesfully stdin")
	}()

	ioReader := os.Stdin
	bufioReader := bufio.NewReader(ioReader)

	readTrigger := make(chan struct{}, 1)
	readTrigger <- struct{}{}

	for {
		select {
		case <-n.closeInboundChan:
			os.Stdin.Close()
			return
		case <-readTrigger:
			go func() {
				// This goroutine gets leaked on closing the service

				in, _, err := bufioReader.ReadLine()
				if err != nil {
					slog.Error("error reading", "error", err)
					return
				}
				slog.Info("received", "in", string(in))

				msg := Message{}
				_ = json.Unmarshal(in, &msg)
				n.inboundChan <- msg

				readTrigger <- struct{}{}
			}()
		}
	}
}

func (n *Node) stdout() {
	defer func() {
		n.closeOutboundChan <- struct{}{}
		close(n.outboundChan)
		slog.Info("closed succesfully stdout")
	}()

	for {
		slog.Info("next handling")
		select {
		case <-n.closeOutboundChan:
			return

		case msg := <-n.outboundChan:
			output, err := json.Marshal(msg)
			if err != nil {
				slog.Error("error marshalling output", "error", err)
				continue
			}
			output = append(output, '\n')

			fmt.Fprint(os.Stdout, string(output))

			slog.Info("sent output", "data", string(output), "n", n)
		}
	}
}

func (n *Node) run() {
	for {
		select {
		case <-n.ctx.Done():
			return
		case msg := <-n.inboundChan:
			slog.Info("handling message", msg)
			err := n.handleMessage(msg)
			if err != nil {
				slog.Error("error handling message", "error", err)
			}
		}
	}
}
