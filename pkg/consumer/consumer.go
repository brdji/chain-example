package consumer

import "github.com/brdji/chain-listener/pkg/message"

type Consumer interface {
	GetQueueName() string
	ListenForMessages()
	ConsumeMessage(*message.DataMessage)
}
