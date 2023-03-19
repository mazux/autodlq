package autodlq

import (
	"github.com/streadway/amqp"
)

type ConsumeConfig interface {
	QueueName() string
	MaxRetryCount() uint64
}

type Consumer interface {
	OnConsuming(d amqp.Delivery) error
	OnConsumed(d amqp.Delivery) error
	OnRetry(d amqp.Delivery) error
	OnMaxRetry(d amqp.Delivery) error
}

func ConsumeQ(channel *amqp.Channel, consumer Consumer, config ConsumeConfig, errChan chan<- error) error {
	messages, err := channel.Consume(
		config.QueueName(), // queue
		"",                 // consumer
		false,              // auto-ack
		false,              // exclusive
		false,              // no-local
		false,              // no-wait
		nil,                // args
	)
	if err != nil {
		return err
	}

	go func() {
		for d := range messages {
			if err := consumer.OnConsuming(d); err != nil {
				if errChan != nil {
					errChan <- err
					errChan <- retryDelivery(d, consumer, config)
				} else {
					_ = retryDelivery(d, consumer, config)
				}
				d.MessageCount++
			} else {
				if errChan != nil {
					errChan <- consumer.OnConsumed(d)
				} else {
					_ = consumer.OnConsumed(d)
				}
			}
		}
	}()

	return nil
}

func retryDelivery(d amqp.Delivery, consumer Consumer, config ConsumeConfig) error {
	xDeath, exists := d.Headers["x-death"].([]interface{})
	if exists {
		c := xDeath[0].(amqp.Table)["count"].(int64)
		if uint64(c) >= config.MaxRetryCount() {
			return consumer.OnMaxRetry(d)
		}
	}

	return consumer.OnRetry(d)
}
