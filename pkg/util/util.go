package util

import "log"

// Function that fails if an error has been returned by RabbitMq
func FailOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}
