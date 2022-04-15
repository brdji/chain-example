package main

import (
	"fmt"
	"log"
	"math/big"
	"net/http"

	"github.com/brdji/chain-example/pkg/chain"
	"github.com/brdji/chain-example/pkg/consumer"
	"github.com/brdji/chain-example/pkg/env"
	"github.com/brdji/chain-example/pkg/listener"
	"github.com/brdji/chain-example/pkg/processor"
	"github.com/brdji/chain-example/pkg/producer"
	"github.com/brdji/chain-example/pkg/rabbit"
	"github.com/brdji/chain-example/pkg/redis"
)

// TODO: store this in redis?
var latestBlock big.Int = *big.NewInt(0)

// Initializes the services in use by the server: blockchain, rabbitMq, redis, etc.
func initialize() {
	serverEnv := env.GetEnv()

	contractAddress := serverEnv.ContractAddress
	chainNodeUrl := serverEnv.ChainNodeUrl
	chain.InitializeContractInstance(contractAddress, chainNodeUrl)

	// TODO: improve wait: redis and rabbit should ping their connections
	log.Printf("Waiting for redis and rabbitMq connections to initialize")
	redis.InitClient()
	rabbit.GetConnection()

	log.Printf("Server initialization complete")
}

// Dummy function performing "work" when a new block changes
func doBlockChangeWork(block big.Int) {
	log.Println("New block #", block.String())
}

func awaitBlockChanges(notifyChan <-chan *big.Int) {
	for blockNum := range notifyChan {
		if blockNum.Cmp(&latestBlock) > 0 {
			latestBlock = *blockNum
			doBlockChangeWork(latestBlock)
		}
	}
}

func startServer() {
	serverEnv := env.GetEnv()

	initialize()

	notifyChan := chain.ListenForBlockChanges(chain.ChainClient)
	go awaitBlockChanges(notifyChan)

	testProcessor()

	port := serverEnv.Port
	log.Printf("Server started and listening at port %s\n", port)

	listenPort := fmt.Sprintf(":%s", port)
	log.Fatal(http.ListenAndServe(listenPort, nil))
}

func testProcessor() {
	// test consumer and producer
	p1 := &producer.RabbitProducer{
		IdGen:     0,
		QueueName: "TestQueue",
	}
	l1 := &listener.RabbitListener{QueueName: "TestQueue"}
	c1 := &consumer.DummyConsumer{}

	processor := &processor.Processor{
		Producer: p1,
		Listener: l1,
		Consumer: c1,
	}

	p1.ProduceMessage("DUMMY PRODUCED DATA 1")
	p1.ProduceMessage("DUMMY PRODUCED DATA 2")
	p1.ProduceMessage("DUMMY PRODUCED DATA 3")
	p1.ProduceMessage("DUMMY PRODUCED DATA 4")
	p1.ProduceMessage("DUMMY PRODUCED DATA 5")

	processor.StartProcessor()
}

func main() {
	startServer()
}
