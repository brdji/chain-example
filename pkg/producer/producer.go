package producer

type Producer interface {
	// Produces a new message from the given string data, and pushes it to queue
	ProduceMessage(string)
}
