package processor

import (
	"github.com/brdji/chain-example/pkg/consumer"
	"github.com/brdji/chain-example/pkg/listener"
	"github.com/brdji/chain-example/pkg/message"
	"github.com/brdji/chain-example/pkg/producer"
)

type Processor struct {
	Producer producer.Producer
	Listener listener.Listener
	Consumer consumer.Consumer
}

func (proc *Processor) StartProcessor() {
	listenerChan := make(chan *message.DataMessage)
	// defer close(listenerChan)

	// send received messages to the consumer
	forever := make(chan bool)

	proc.Listener.SetNotifyChannel(listenerChan)

	go func() {
		for msg := range listenerChan {
			if msg != nil {
				proc.Consumer.ConsumeMessage(msg)
			}
		}
	}()

	go func() {
		proc.Listener.ListenForMessages()
	}()

	<-forever
}
