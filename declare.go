package autodlq

import (
	"github.com/streadway/amqp"
)

type DeclareConfig interface {
	QueueName() string
	DLQName() string
	RoutingKeyName() string
	ExchangeName() string
	DLXName() string
	MsgTTL() uint64
	SupportDlq() bool
}

func DeclareQ(channel *amqp.Channel, c DeclareConfig) error {
	err := declareQueueWithBind(channel, c)
	if err != nil {
		return err
	}

	if c.SupportDlq() {
		err = declareDLWithBind(channel, c)
		if err != nil {
			return err
		}
	}

	return nil
}

func declareQueueWithBind(channel *amqp.Channel, c DeclareConfig) error {
	err := declareExchange(channel, c.ExchangeName())
	if err != nil {
		return err
	}

	err = declareQueueWithDLX(channel, c.QueueName(), c.DLXName(), 0)
	if err != nil {
		return err
	}

	err = bindQueueToExchange(channel, c.QueueName(), c.RoutingKeyName(), c.ExchangeName())
	if err != nil {
		return err
	}

	return nil
}

func declareDLWithBind(channel *amqp.Channel, c DeclareConfig) error {
	err := declareExchange(channel, c.DLXName())
	if err != nil {
		return err
	}

	err = declareQueueWithDLX(channel, c.DLQName(), c.ExchangeName(), c.MsgTTL())
	if err != nil {
		return err
	}

	err = bindQueueToExchange(channel, c.DLQName(), c.RoutingKeyName(), c.DLXName())
	if err != nil {
		return err
	}

	return nil
}

func declareQueueWithDLX(channel *amqp.Channel, queueName, dlxName string, msgTTL uint64) error {
	args := amqp.Table{}
	if dlxName != "" {
		args["x-dead-letter-exchange"] = dlxName
	}
	if msgTTL > 0 {
		args["x-message-ttl"] = int64(msgTTL)
	}

	_, err := channel.QueueDeclare(
		queueName, // queue name
		true,      // durable
		false,     // auto delete
		false,     // exclusive
		false,     // no wait
		args,      // arguments
	)

	return err
}

func declareExchange(channel *amqp.Channel, exchangeName string) error {
	return channel.ExchangeDeclare(
		exchangeName, // name
		"topic",      // type
		true,         // durable
		false,        // auto-deleted
		false,        // internal
		false,        // no-wait
		nil,          // arguments
	)
}

func bindQueueToExchange(channel *amqp.Channel, queueName, routingKey, exchangeName string) error {
	return channel.QueueBind(
		queueName,
		routingKey,
		exchangeName,
		false,
		nil,
	)
}
