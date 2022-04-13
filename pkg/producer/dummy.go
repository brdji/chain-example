package producer

import (
	"encoding/json"
	"fmt"

	"github.com/brdji/chain-example/pkg/message"
	"github.com/brdji/chain-example/pkg/rabbit"
	"github.com/brdji/chain-example/pkg/util"
)

type DummyProducer struct {
	IdGen int32
}

func (prod *DummyProducer) GetQueueName() string {
	return "dummy-messages"
}

func (prod *DummyProducer) ProduceMessage(data string) {
	prod.IdGen++
	msg := &message.DataMessage{
		Id:   fmt.Sprint(prod.IdGen),
		Data: data,
	}

	jsonData, err := json.Marshal(msg)
	util.FailOnError(err, "Error producing message")

	rabbit.PublishMessage(prod.GetQueueName(), jsonData)
}
