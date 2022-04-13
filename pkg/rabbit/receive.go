package rabbit

import (
	"log"

	"github.com/brdji/chain-listener/pkg/util"

	"github.com/streadway/amqp"
)

func TestReceive() {
	conn := GetConnection()

	ch, err := conn.Channel()
	util.FailOnError(err, "Error opening rabbitMq channel")

	defer ch.Close()

	msgs, err := ch.Consume(
		TestQueueName, // queue
		"",            // zconsumer
		true,          // auto ack
		false,         // exclusive
		false,         // no local
		false,         // no wait
		nil,           // args
	)
	util.FailOnError(err, "Failed to register a consumer")
	forever := make(chan bool)
	go func() {
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)
			forever <- true
		}
	}()

	log.Printf("Waiting for messages forever")
	<-forever
}

func GetQueueChannel(queueName string) <-chan amqp.Delivery {
	conn := GetConnection()

	ch, err := conn.Channel()
	util.FailOnError(err, "Error opening rabbitMq channel")

	deliveryChan, err := ch.Consume(
		TestQueueName, // queue
		"",            // zconsumer
		true,          // auto ack
		false,         // exclusive
		false,         // no local
		false,         // no wait
		nil,           // args
	)
	util.FailOnError(err, "Failed to register a consumer")

	return deliveryChan
}
