package listener

import "github.com/brdji/chain-example/pkg/message"

type Listener interface {
	ListenForMessages()
	SetNotifyChannel(chan<- *message.DataMessage)
}
