package consumer

import (
	"encoding/json"
	"log"

	"github.com/brdji/chain-listener/pkg/message"
	"github.com/brdji/chain-listener/pkg/rabbit"
	"github.com/brdji/chain-listener/pkg/util"
)

type DummyConsumer struct {
}

func (consumer *DummyConsumer) GetQueueName() string {
	return "dummy-messages"
}

func (consumer *DummyConsumer) ListenForMessages() {

	listenChan := rabbit.GetQueueChannel(consumer.GetQueueName())

	forever := make(chan bool)

	go func() {
		for chanMsg := range listenChan {
			var msg *message.DataMessage
			err := json.Unmarshal(chanMsg.Body, &msg)
			util.FailOnError(err, "Error decoding chan message")

			consumer.ConsumeMessage(msg)
		}
	}()

	<-forever
}

func (consumer *DummyConsumer) ConsumeMessage(msg *message.DataMessage) {
	log.Println("Consumed message #", msg.Id, ". Data: ", msg.Data)
}
