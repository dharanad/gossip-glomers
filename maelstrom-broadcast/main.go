package main

import (
	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
	"log"
)

func main() {
	n := maelstrom.NewNode()
	n.Handle("broadcast", NewBroadcastHandler(n))
	n.Handle("read", NewReadHandler(n))
	n.Handle("topology", NewTopologyHandler(n))
	if err := n.Run(); err != nil {
		log.Fatal(err)
	}
}

func NewBroadcastHandler(n *maelstrom.Node) maelstrom.HandlerFunc {
	return func(msg maelstrom.Message) error {
		return n.Reply(msg, map[string]any{})
	}
}

func NewReadHandler(n *maelstrom.Node) maelstrom.HandlerFunc {
	return func(msg maelstrom.Message) error {
		return n.Reply(msg, map[string]any{})
	}
}

func NewTopologyHandler(n *maelstrom.Node) maelstrom.HandlerFunc {
	return func(msg maelstrom.Message) error {
		return n.Reply(msg, map[string]any{})
	}
}
