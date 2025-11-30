package rabbitmq

import "github.com/rabbitmq/amqp091-go"

type Consumer struct {
	Conn *amqp091.Connection
	Chan *amqp091.Channel
	Msgs <-chan amqp091.Delivery
}

// NewConsumer connects and create a consumer for a queue
func NewConsumer(url, queueName string) (*Consumer, error) {
	conn, err := amqp091.Dial(url)
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	// Ensure queue exists
	_, err = ch.QueueDeclare(
		queueName,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, err
	}

	msgs, err := ch.Consume(
		queueName,
		"",
		false, // manual ACK
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, err
	}

	return &Consumer{
		Conn: conn,
		Chan: ch,
		Msgs: msgs,
	}, nil
}

func (c *Consumer) Close() {
	_ = c.Chan.Close()
	_ = c.Conn.Close()
}
