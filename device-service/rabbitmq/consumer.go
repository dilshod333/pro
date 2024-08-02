package rabbitmq

import (
	"context"
	"device-service/pkg"
	"device-service/protos/gendevice"
	"encoding/json"
	"log"

	"github.com/streadway/amqp"
)

type  Serverr struct{
	S *pkg.Server
}

func Justt() *Serverr{
	client  := pkg.Server{}
	return &Serverr{S: &client}
}


func ConsumeOrders() {
	conn, err := amqp.Dial("amqp://guest:guest@message-broker:5672/")
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Consumer ichiga tushdi...")

	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatal(err)
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
		log.Fatal(err)
	}

	q, err := ch.QueueDeclare(
		"",
		false,
		false,
		true,
		false,
		nil,
	)
	if err != nil {
		log.Fatal(err)
	}

	err = ch.QueueBind(
		q.Name,
		"",
		"command",
		false,
		nil,
	)
	if err != nil {
		log.Fatal(err)
	}

	msgs, err := ch.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatal(err)
	}

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			var command *gendevice.DeviceControlReq
			err := json.Unmarshal(d.Body, &command)
			if err != nil {
				log.Printf("Error decoding JSON: %s", err)
				continue
			}
			log.Println("IChida consumerrr ishladiiii")


			err = Justt().S.UpdateControl(context.Background(), command)
			if err != nil {
				log.Println("Xatolik updatecontrollda...", err)
				continue 
			}
			

			log.Printf("Buyurtma olindi va qayta ishlanmoqda: %+v", command)
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}
