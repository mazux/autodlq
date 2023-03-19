package main

type config struct {
	queue      string `env:"QUEUE_NAME"`
	maxRetries uint64 `env:"MAX_RETRIES"`
}

func (c *config) QueueName() string {
	return c.queue
}

func (c *config) MaxRetryCount() uint64 {
	return c.maxRetries
}

func NewConfig(queue string, maxRetries uint64) *config {
	return &config{
		queue:      queue,
		maxRetries: maxRetries,
	}
}
