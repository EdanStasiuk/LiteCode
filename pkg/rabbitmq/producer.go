package rabbitmq

import (
	"fmt"
	"log"

	"github.com/rabbitmq/amqp091-go"
)

var producerConn *amqp091.Connection
var producerChan *amqp091.Channel
var producerQueue amqp091.Queue

// InitProducer initializes a RabbitMQ producer
func InitProducer(url string, queueName string) error {
	var err error

	producerConn, err = amqp091.Dial(url)
	if err != nil {
		return err
	}

	producerChan, err = producerConn.Channel()
	if err != nil {
		return err
	}

	producerQueue, err = producerChan.QueueDeclare(
		queueName,
		true,  // durable
		false, // autoDelete
		false, // exclusive
		false, // noWait
		nil,   // args
	)
	if err != nil {
		return err
	}

	log.Printf("RabbitMQ producer initialized for queue: %s", producerQueue.Name)

	return nil
}

func ProduceMessage(body []byte) error {
	if producerChan == nil {
		fmt.Printf("Producer channel is nil\n")
		return fmt.Errorf("producer channel is nil; have you called InitProducer() successfully?")
	}
	return producerChan.Publish(
		"",                 // exchange
		producerQueue.Name, // routing key
		false,              // mandatory
		false,              // immediate
		amqp091.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
}

func CloseProducer() {
	if producerChan != nil {
		_ = producerChan.Close()
	}
	if producerConn != nil {
		_ = producerConn.Close()
	}
}
