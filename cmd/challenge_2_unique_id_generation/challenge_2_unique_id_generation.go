package main

import (
	"fmt"
	"log"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

func main() {
	node := maelstrom.NewNode()
	count := int64(0)
	node.Handle("generate", func(msg maelstrom.Message) error {
		body := map[string]any{
			"type": "generate_ok",
			"id":   fmt.Sprintf("%s/%d", node.ID(), count),
		}
		count++
		return node.Reply(msg, body)
	})
	if err := node.Run(); err != nil {
		log.Fatal(err)
	}
}
