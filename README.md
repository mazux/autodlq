# Auto DLQ

This is a package I built while working on one of my article. I was trying things around and explore rabbitMQ.
It depends on [streadway/amqp](https://github.com/streadway/amqp).

## autodlq.DeclareQ

It helps to declare a queue with its DLX in case you passed a config implementation with `SupportDlq() bool` that returns `true`.

## autodlq.ConsumeQ

It helps decouple your infrastructure code from application. You are not supposed to worry about infrastructure things. You just need to have `Consumer` implementation with the desired behaviour you need on:

- `OnConsuming(d amqp.Delivery) error` What functions should be invoked while handling your message.
- `OnConsumed(d amqp.Delivery) error` What should happen after successfully consuming the message. Maybe some useful logs. 
- `OnRetry(d amqp.Delivery) error` How to retry the message.
- `OnMaxRetry(d amqp.Delivery) error` What should happen when we reached the max retry.
