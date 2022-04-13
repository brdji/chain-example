package rabbit

import (
	"log"
	"time"

	env "github.com/brdji/chain-listener/pkg/env"
	"github.com/brdji/chain-listener/pkg/util"

	"github.com/streadway/amqp"
)

var rabbitConnection *amqp.Connection

const TestQueueName string = "testQueue"

func initConnection(conn *amqp.Connection) {
	ch, err := conn.Channel()
	util.FailOnError(err, "Error opening rabbitMq channel")

	defer ch.Close()

	// declare the queues we need
	_, err = ch.QueueDeclare(
		TestQueueName, // name
		false,         // durable
		false,         // autodelete
		false,         // exclusive
		false,         // no wait
		nil,           // arguments
	)
	util.FailOnError(err, "Error declaring queue")
}

// Returns the connection to RabbitMq. If the connection is not established, will attempt to
// establish it, with a specified timeout period
func GetConnection() *amqp.Connection {
	if rabbitConnection == nil {
		var err error
		for i := 10; i > 0; i-- {
			rabbitConnection, err = amqp.Dial(env.GetEnv().RabbitBrokerUrl)
			if err == nil {
				break
			} else {
				time.Sleep(time.Second * 2)
			}
		}

		// if err still exists, then we failed to connect to rabbit
		util.FailOnError(err, "Error connecting to rabbitMq")
		initConnection(rabbitConnection)
		log.Println("Connected to rabbitMQ")
	}
	return rabbitConnection
}
