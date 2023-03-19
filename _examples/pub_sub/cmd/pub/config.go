package main

import "fmt"

type config struct {
	queue    string `env:"QUEUE_NAME"`
	routingK string `env:"ROUTING_KEY"`
	exchange string `env:"EXCHANGE_NAME"`
	ttl      uint64 `env:"EXCHANGE_NAME"`
}

func (c *config) QueueName() string {
	return c.queue
}

func (c *config) DLQName() string {
	return fmt.Sprintf("%s_DLQ", c.queue)
}

func (c *config) RoutingKeyName() string {
	return c.routingK
}

func (c *config) ExchangeName() string {
	return c.exchange
}

func (c *config) DLXName() string {
	return fmt.Sprintf("%s_DLX", c.exchange)
}

func (c *config) MsgTTL() uint64 {
	return c.ttl
}

func (c *config) SupportDlq() bool {
	return c.ttl > 0
}

func NewConfig(queue, routingK, exchange string, ttl uint64) *config {
	return &config{
		queue:    queue,
		routingK: routingK,
		exchange: exchange,
		ttl:      ttl,
	}
}
