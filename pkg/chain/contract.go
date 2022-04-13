package chain

import (
	"log"

	"github.com/brdji/chain-example/pkg/contract"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

var ContractInstance *contract.Contract
var ChainClient *ethclient.Client

func InitializeContractInstance(contractHexAddress string, clientNodeUrl string) {
	contractAddress := common.HexToAddress(contractHexAddress)

	// TODO how do I prevent shadowing the ChainClient by avoiding the := assignment to new variables, below?
	var err error
	ChainClient, err = ethclient.Dial(clientNodeUrl)
	if err != nil {
		log.Fatal(err)
	}

	ContractInstance, err = contract.NewContract(contractAddress, ChainClient)

	if err != nil {
		log.Fatal(err)
	}

	log.Println("Contract instance initialized!")
}
