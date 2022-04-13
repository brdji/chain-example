package env

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Env struct {
	Port            string
	ContractAddress string
	ChainNodeUrl    string
	RabbitBrokerUrl string
	RedisUrl        string
}

var listenerEnv *Env

func InitEnv() {
	// if we're running locally, load the env from file
	_, envSpecified := os.LookupEnv("env_type")
	if !envSpecified {
		log.Println("Env unspecified: loading from file")
		err := godotenv.Load("run.env")
		if err != nil {
			log.Fatal("Error loading env : ", err)
		}
	}

	listenerEnv = &Env{}
	listenerEnv.Port = os.Getenv("port")
	listenerEnv.ContractAddress = os.Getenv("contract_address")
	listenerEnv.ChainNodeUrl = os.Getenv("chain_node_url")
	listenerEnv.RabbitBrokerUrl = os.Getenv("rabbit_broker_url")
	listenerEnv.RedisUrl = os.Getenv("redis_url")
}

func GetEnv() *Env {
	if listenerEnv == nil {
		InitEnv()
	}
	return listenerEnv
}
