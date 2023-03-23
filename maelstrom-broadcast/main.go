package main

import (
	"encoding/json"
	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
	"log"
	"sync"
)

func main() {
	n := maelstrom.NewNode()
	n.NodeIDs()
	svc := NewService(n)
	n.Handle("broadcast", svc.BroadcastHandler())
	n.Handle("read", svc.ReadHandler())
	n.Handle("topology", svc.TopologyHandler())
	if err := n.Run(); err != nil {
		log.Fatal(err)
	}
}

type Service struct {
	n             *maelstrom.Node
	ids           []int
	idsMutex      sync.RWMutex
	topology      map[string][]string
	topologyMutex sync.RWMutex
}

func NewService(n *maelstrom.Node) *Service {
	return &Service{n: n}
}

func (s *Service) BroadcastHandler() maelstrom.HandlerFunc {
	return func(msg maelstrom.Message) error {
		var body map[string]any
		if err := json.Unmarshal(msg.Body, &body); err != nil {
			return err
		}
		// Ref: https://go.dev/ref/spec#Type_assertions
		message := body["message"].(float64)
		s.idsMutex.Lock()
		s.ids = append(s.ids, int(message))
		s.idsMutex.Unlock()
		return s.n.Reply(msg, map[string]any{
			"type": "broadcast_ok",
		})
	}
}

func (s *Service) ReadHandler() maelstrom.HandlerFunc {
	return func(msg maelstrom.Message) error {
		var body map[string]any
		if err := json.Unmarshal(msg.Body, &body); err != nil {
			return err
		}
		s.idsMutex.RLock()
		ids := make([]int, len(s.ids))
		for idx, id := range s.ids {
			ids[idx] = id
		}
		s.idsMutex.RUnlock()
		return s.n.Reply(msg, map[string]any{
			"type":     "read_ok",
			"messages": ids,
		})
	}
}

func (s *Service) TopologyHandler() maelstrom.HandlerFunc {
	return func(msg maelstrom.Message) error {
		var body struct {
			Topology map[string][]string `json:"topology"`
		}
		if err := json.Unmarshal(msg.Body, &body); err != nil {
			return err
		}
		s.topologyMutex.Lock()
		s.topology = body.Topology
		s.topologyMutex.Unlock()
		return s.n.Reply(msg, map[string]any{
			"type": "topology_ok",
		})
	}
}
