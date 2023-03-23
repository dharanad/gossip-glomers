package main

import (
	"encoding/json"
	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
	"log"
	"sync/atomic"
)

func main() {
	n := maelstrom.NewNode()
	var idGenerator IdGenerator = &UniqueIdGeneratorService{}
	n.Handle("generate", NewIdGeneratorHandler(idGenerator, n))
	if err := n.Run(); err != nil {
		log.Fatal(err)
	}
}

func NewIdGeneratorHandler(idGenerator IdGenerator, n *maelstrom.Node) maelstrom.HandlerFunc {
	return func(msg maelstrom.Message) error {
		var body map[string]any
		if err := json.Unmarshal(msg.Body, &body); err != nil {
			return err
		}
		body["type"] = "generate_ok"
		body["id"] = idGenerator.Generate()
		return n.Reply(msg, body)
	}
}

type IdGenerator interface {
	Generate() int64
}

type UniqueIdGeneratorService struct {
	count int64
}

func (u *UniqueIdGeneratorService) Generate() int64 {
	return atomic.AddInt64(&u.count, 1)
}
