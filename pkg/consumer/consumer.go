package consumer

import "github.com/brdji/chain-example/pkg/message"

type Consumer interface {
	ConsumeMessage(*message.DataMessage)
}
