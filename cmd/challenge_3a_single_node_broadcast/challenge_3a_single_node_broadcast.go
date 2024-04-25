package main

import (
	"encoding/json"
	"log"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

type BroadcastRequestBody struct {
	Message float64 `json:"message"`
}

type TopologyRequestBody struct {
	Topology map[string][]string `json:"topology"`
}

func main() {
	node := maelstrom.NewNode()
	neighbours := make([]string, 0)
	messages := make([]int, 0)
	node.Handle("broadcast", func(msg maelstrom.Message) error {
		var body BroadcastRequestBody
		if err := json.Unmarshal(msg.Body, &body); err != nil {
			return err
		}
		messages = append(messages, int(body.Message))
		for _, id := range neighbours {
			if id == node.ID() {
				continue
			}
			node.Send(id, msg)
		}
		return node.Reply(msg, map[string]string{"type": "broadcast_ok"})
	})
	node.Handle("read", func(msg maelstrom.Message) error {
		return node.Reply(msg, map[string]any{
			"type":     "read_ok",
			"messages": messages,
		})
	})
	node.Handle("topology", func(msg maelstrom.Message) error {
		var body TopologyRequestBody
		if err := json.Unmarshal(msg.Body, &body); err != nil {
			return err
		}
		neighbours = body.Topology[node.ID()]
		return node.Reply(msg, map[string]string{"type": "topology_ok"})
	})
	if err := node.Run(); err != nil {
		log.Fatal(err)
	}
}
