package message

// A single message created by a producer, handled by a consumer
type DataMessage struct {
	Id   string
	Data string // represents any type of data, e.g. JSON encoded data
}
