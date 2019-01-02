package benchmark

type MessageQueue interface {
	// MessageReceiver interface
	MessageHandler() *MessageHandler
	Receive()
	// MessageSender interface
	Send([]byte, bool, func(string, error))
	// Common interface
	Teardown()
}
