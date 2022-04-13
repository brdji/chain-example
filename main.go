package main

import (
	"fmt"
	"log"
	"math/big"
	"net/http"

	"github.com/brdji/chain-example/pkg/chain"
	"github.com/brdji/chain-example/pkg/consumer"
	"github.com/brdji/chain-example/pkg/env"
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

	// test consumer and producer
	p1 := &producer.DummyProducer{
		IdGen: 0,
	}
	p2 := &producer.DummyProducer{
		IdGen: 10,
	}
	consumer := &consumer.DummyConsumer{}

	p1.ProduceMessage("DUMMY PRODUCED DATA")
	p2.ProduceMessage("DUMMY PRODUCED DATA 2")
	p2.ProduceMessage("DUMMY PRODUCED DATA 3")
	p1.ProduceMessage("DUMMY PRODUCED DATA 4")
	go func() { consumer.ListenForMessages() }()

	port := serverEnv.Port
	log.Printf("Server started and listening at port %s\n", port)

	listenPort := fmt.Sprintf(":%s", port)
	log.Fatal(http.ListenAndServe(listenPort, nil))
}

func main() {
	startServer()
}
