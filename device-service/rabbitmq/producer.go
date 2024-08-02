package rabbitmq

import (
	"log"

	"github.com/streadway/amqp"
)

func PublishOrder(body []byte) error {
	conn, err := amqp.Dial("amqp://guest:guest@message-broker:5672/")
	if err != nil {
		return err
	}
	log.Println("Publish ichiga tushdi...")

	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	err = ch.ExchangeDeclare(
		"command",
		"fanout",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	err = ch.Publish(
		"command",
		"",
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		})
	if err != nil {
		return err
	}

	return nil
}
