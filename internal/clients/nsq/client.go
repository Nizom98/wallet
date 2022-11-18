package nsq

import (
	"fmt"
	"github.com/nsqio/go-nsq"
)

// Client клиент для взаимодействия c nsq.
type Client struct {
	producer *nsq.Producer
	topic string
}

// NewClient конструктор клиента nsq.
// topic - топик, в который будет производится запись.
// target - эндпоинт для подключения к nsq.
func NewClient(topic, target string) (*Client, error) {
	if topic == "" {
		return nil, fmt.Errorf("empty topic")
	}
	if target == "" {
		return nil, fmt.Errorf("empty target")
	}

	cfg := nsq.NewConfig()
	producer, err := nsq.NewProducer(target, cfg)
	if err != nil {
		return nil, fmt.Errorf("cannot create producer for %s: %w", topic, err)
	}

	client := &Client{
		producer: producer,
		topic:    topic,
	}

	return client, nil
}

func (c *Client) Stop() {
	c.producer.Stop()
}

// Write отправка сообщения в nsq.
// Топик в который будет отправляться, определяется в конструкторе.
func (c *Client) Write(data []byte) error {
	err := c.producer.Publish(c.topic, data)
	if err != nil {
		return fmt.Errorf("cannot write to %s: %w", c.topic, err)
	}

	return nil
}
