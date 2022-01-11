package main

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/marcoshuck/nats-demo/pkg/models"
	"github.com/nats-io/nats.go"
	"log"
	"time"
)

func main() {
	conn, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		log.Fatalln("Failed to connect to NATS cluster:", err)
	}
	js, err := conn.JetStream()
	if err != nil {
		log.Fatalln("Failed to create NATS JetStream context:", err)
	}

	for i := 0; i < 10; i++ {
		order := models.Order{
			UUID:      uuid.New(),
			Status:    "received",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		var b []byte
		if b, err = json.Marshal(order); err != nil {
			log.Println("Error marshaling order:", order.UUID)
			continue
		}

		_, err = js.Publish("ORDERS.received", b)
		if err != nil {
			log.Println("Error publishing order:", order.UUID, ". Error:", err)
			continue
		}
	}
}
