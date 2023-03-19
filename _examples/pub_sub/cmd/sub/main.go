package main

import (
	"fmt"
	"log"
	"os"

	"github.com/mazux/autodlq"
	"github.com/streadway/amqp"
)

var (
	queueName     = "queue_name.event.action"
	amqpServerURL = "amqp://guest:guest@localhost:5672/"
)

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

	err = channelRabbitMQ.Qos(1, 0, false)
	if err != nil {
		fmt.Println("unable to set QoS on a channel. Error: ", err)
		os.Exit(1)
	}

	var forever chan bool

	errChan := make(chan error)
	go func() {
		for err = range errChan {
			if err != nil {
				log.Println("Error from consumer. Error:", err)
			}
		}
	}()

	config := NewConfig(queueName, 3)
	consumer := NewConsumer()
	err = autodlq.ConsumeQ(channelRabbitMQ, consumer, config, errChan)
	if err != nil {
		fmt.Println("unable to consume from a channel. Error: ", err)
		os.Exit(1)
	}

	<-forever
}
