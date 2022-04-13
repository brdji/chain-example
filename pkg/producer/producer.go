package producer

type Producer interface {
	// Returns the queue name where the producer will push messages
	GetQueueName() string
	// Produces a new message from the given string data, and pushes it to queue
	ProduceMessage(string)
}
