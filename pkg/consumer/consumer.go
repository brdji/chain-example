package consumer

import "github.com/brdji/chain-example/pkg/message"

type Consumer interface {
	GetQueueName() string
	ListenForMessages()
	ConsumeMessage(*message.DataMessage)
}
