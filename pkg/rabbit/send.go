package rabbit

import (
	"github.com/brdji/chain-listener/pkg/util"

	"github.com/streadway/amqp"
)

func PublishMessage(queueName string, data []byte) {
	conn := GetConnection()

	ch, err := conn.Channel()
	util.FailOnError(err, "Error opening rabbitMq channel")

	defer ch.Close()

	err = ch.Publish(
		"",            // exchange
		TestQueueName, // routing key
		false,         // mandatory
		false,         // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        data,
		},
	)

	util.FailOnError(err, "Error publishing message")
}
