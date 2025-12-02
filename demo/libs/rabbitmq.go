package libs

import (
	"log"

	"github.com/streadway/amqp"
)

type RabbitMQClient struct {
	conn *amqp.Connection
	ch   *amqp.Channel
}

func NewRabbitMQClient(url string, prefetchCount int) *RabbitMQClient {
	if prefetchCount <= 0 {
		prefetchCount = 100
	}

	conn, err := amqp.Dial(url)
	if err != nil {
		log.Panic(err)
	}

	ch, err := conn.Channel()
	if err != nil {
		log.Panic(err)
	}

	ch.Qos(prefetchCount, 0, false)

	return &RabbitMQClient{
		conn: conn,
		ch:   ch,
	}
}

func (c *RabbitMQClient) Close() {
	c.ch.Close()
	c.conn.Close()
}

func (c *RabbitMQClient) QueueDeclare(name string) amqp.Queue {
	q, err := c.ch.QueueDeclare(
		name,
		false,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		log.Panic(err)
	}
	return q
}

func (c *RabbitMQClient) Consume(queue string, consumer string) <-chan amqp.Delivery {
	msgs, err := c.ch.Consume(
		queue,
		consumer,
		true,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		log.Panic(err)
	}
	return msgs
}

func (c *RabbitMQClient) Channel() *amqp.Channel {
	return c.ch
}
