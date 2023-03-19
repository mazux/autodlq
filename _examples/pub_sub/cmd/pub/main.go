package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/mazux/autodlq"
	"github.com/streadway/amqp"
)

var (
	exchange      = "my-exchange"
	queueName     = "queue_name.event.action"
	amqpServerURL = "amqp://guest:guest@localhost:5672/"
)

type ResourceCreatedEvent struct {
	Id int `json:"id"`
}

func main() {
	// Create a new RabbitMQ connection.
	connectRabbitMQ, err := amqp.Dial(amqpServerURL)
	if err != nil {
		fmt.Println("unable to connect to rabbitMQ. Error: ", err)
		os.Exit(1)
	}
	defer connectRabbitMQ.Close()

	channelRabbitMQ, err := connectRabbitMQ.Channel()
	if err != nil {
		fmt.Println("unable to open a channel. Error: ", err)
		os.Exit(1)
	}
	defer channelRabbitMQ.Close()

	config := NewConfig(queueName, "rout", exchange, 5000)
	err = autodlq.DeclareQ(channelRabbitMQ, config)
	if err != nil {
		fmt.Println("unable to auto declare queue. Error: ", err)
		os.Exit(1)
	}

	msg := ResourceCreatedEvent{
		Id: 10,
	}

	blob, err := json.Marshal(msg)
	if err != nil {
		fmt.Println("unable to marshal the message. Error: ", err)
		os.Exit(1)
	}

	message := amqp.Publishing{
		//ContentType: "application/json",
		Body: blob,
		//Type:        "message type",
		//MessageId:   "uuid",
		//Expiration: "9000",
	}

	err = channelRabbitMQ.Publish(
		exchange, // exchange
		"rout",   // queue name
		false,    // mandatory
		false,    // immediate
		message,  // message to publish
	)

	if err != nil {
		fmt.Println("unable to publish message to queue. Error: ", err)
		os.Exit(1)
	}
}
