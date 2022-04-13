package chain

import (
	"context"
	"log"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/ethclient"
)

func getLatestBlockNumber(client *ethclient.Client) *big.Int {
	header, err := client.HeaderByNumber(context.Background(), nil)
	if err != nil {
		log.Fatal(err)
	}
	return header.Number
}

func ListenForBlockChanges(client *ethclient.Client) <-chan *big.Int {
	log.Println("Listening for block changes")

	notifyChan := make(chan *big.Int, 1)

	go func() {
		for {
			blockNum := getLatestBlockNumber(client)
			notifyChan <- blockNum
			// don't spam too much, sleep for a while before asking for block updates
			time.Sleep(10 * time.Second)
		}
	}()
	return notifyChan
}
