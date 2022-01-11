package main

import (
	"encoding/json"
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

	js.AddConsumer("orders", &nats.ConsumerConfig{
		Durable:         "",
		Description:     "",
		DeliverSubject:  "",
		DeliverGroup:    "",
		DeliverPolicy:   0,
		OptStartSeq:     0,
		OptStartTime:    nil,
		AckPolicy:       0,
		AckWait:         0,
		MaxDeliver:      0,
		FilterSubject:   "",
		ReplayPolicy:    0,
		RateLimit:       0,
		SampleFrequency: "",
		MaxWaiting:      0,
		MaxAckPending:   0,
		FlowControl:     false,
		Heartbeat:       0,
		HeadersOnly:     false,
	})
	_, err = js.Subscribe("ORDERS.*", deliverOrder(js), nats.Durable("ORDERS"), nats.MaxDeliver(3))
	if err != nil {
		log.Fatalln("Failed to subscribe to ORDERS.*:", err)
	}

	select {}
}

func deliverOrder(js nats.JetStream) nats.MsgHandler {
	return func(msg *nats.Msg) {
		var order models.Order
		err := json.Unmarshal(msg.Data, &order)
		if err != nil {
			log.Println("Invalid order:", err)
			return
		}

		switch msg.Subject {
		case "ORDERS.received":
			order.Status = "completed"
			order.UpdatedAt = time.Now()
			b, err := json.Marshal(order)
			if err != nil {
				log.Println("Failed to marshal order:", err)
				if err = msg.Nak(); err != nil {
					log.Println("[FATAL] Failed to NAK")
					return
				}
				return
			}
			_, err = js.Publish("ORDERS.completed", b)
			if err != nil {
				log.Println("Failed to publish completed order:", err)
				if err = msg.Nak(); err != nil {
					log.Println("[FATAL] Failed to NAK")
					return
				}
				return
			}
		case "ORDERS.completed":
			log.Println("Order completed:", order.UUID)
			if err = msg.Ack(); err != nil {
				log.Println("[FATAL] Failed to ACK")
				return
			}
		}
	}
}
