package listener

import (
	"encoding/json"

	"github.com/brdji/chain-example/pkg/message"
	"github.com/brdji/chain-example/pkg/rabbit"
	"github.com/brdji/chain-example/pkg/util"
)

// Simple listener that listens for messages on a specified RabbitMq queue
type RabbitListener struct {
	QueueName     string
	NotifyChannel chan<- *message.DataMessage
}

func (listener *RabbitListener) ListenForMessages() {

	listenChan := rabbit.GetQueueChannel(listener.QueueName)

	forever := make(chan bool)

	go func() {
		for chanMsg := range listenChan {
			var msg *message.DataMessage
			err := json.Unmarshal(chanMsg.Body, &msg)
			util.FailOnError(err, "Error decoding chan message")

			listener.NotifyChannel <- msg
		}
	}()

	<-forever
}

func (listener *RabbitListener) SetNotifyChannel(channel chan<- *message.DataMessage) {
	listener.NotifyChannel = channel
}
