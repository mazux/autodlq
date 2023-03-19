package main

import (
	"fmt"
	"log"

	"github.com/streadway/amqp"
)

type consumer struct {
}

func (c *consumer) OnConsuming(d amqp.Delivery) error {
	log.Println("onConsuming...")

	return fmt.Errorf("error while consuming")

	//return d.Ack(false)
}

func (c *consumer) OnConsumed(d amqp.Delivery) error {
	log.Println("onConsumed...")
	log.Println("consumed successfully")
	return nil
}

func (c *consumer) OnRetry(d amqp.Delivery) error {
	log.Println("onRetry...")

	return d.Nack(false, false)
}

func (c *consumer) OnMaxRetry(d amqp.Delivery) error {
	log.Println("onMaxRetry...")

	return d.Ack(false)
}

func NewConsumer() *consumer {
	return &consumer{}
}
