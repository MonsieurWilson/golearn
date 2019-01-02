package benchmark

import (
	"log"
	"sync"
	"time"
)

type MessageReceiver interface {
	MessageHandler() *MessageHandler
	Receive()
	Teardown()
}

type ReceiveEndpoint struct {
	MessageReceiver  MessageReceiver
	NumberOfMessages int
	SizeOfMessage    int
	Async            bool
	Handler          *MessageHandler
}

func NewReceiveEndpoint(tester Tester) *ReceiveEndpoint {
	return &ReceiveEndpoint{
		MessageReceiver:  tester.Receiver,
		NumberOfMessages: tester.MessageCount,
		SizeOfMessage:    tester.MessageSize,
		Async:            tester.Async,
		Handler:          tester.Receiver.MessageHandler(),
	}
}

type MessageHandler interface {
	// Process a received message. Return true if it's the last message, otherwise
	// return false.
	ReceiveMessage([]byte) bool

	// Indicate whether the handler has been marked complete, meaning all messages
	// have been received.
	HasCompleted() bool
}

type ThroughputMessageHandler struct {
	hasStarted       bool
	hasCompleted     bool
	messageCounter   int
	NumberOfMessages int
	SizeOfMessage    int
	started          int64
	stopped          int64
	completionLock   sync.Mutex
}

func (handler *ThroughputMessageHandler) HasCompleted() bool {
	handler.completionLock.Lock()
	defer handler.completionLock.Unlock()
	return handler.hasCompleted
}

// Increment a message counter. If this is the first message, set the started timestamp.
// If it's the last message, set the stopped timestamp and compute the total runtime
// and print it out. Return true if it's the last message, otherwise return false.
func (handler *ThroughputMessageHandler) ReceiveMessage(message []byte) bool {
	if !handler.hasStarted {
		handler.hasStarted = true
		handler.started = time.Now().UnixNano()
	}

	handler.messageCounter++

	if handler.messageCounter == handler.NumberOfMessages {
		handler.stopped = time.Now().UnixNano()
		elapsed := float32(handler.stopped-handler.started) / 1e9
		log.Printf("Received %d messages (size %d) in %.2f s",
			handler.NumberOfMessages,
			handler.SizeOfMessage,
			elapsed)
		log.Printf("Rate %.2f Mbps, %.2f msg/s",
			float32(8*handler.SizeOfMessage*handler.NumberOfMessages)/elapsed/1e6,
			float32(handler.NumberOfMessages)/elapsed)
		handler.completionLock.Lock()
		handler.hasCompleted = true
		handler.completionLock.Unlock()

		return true
	}

	return false
}

func (endpoint ReceiveEndpoint) WaitForCompletion() {
	for {
		if (*endpoint.Handler).HasCompleted() {
			break
		} else {
			time.Sleep(10 * time.Millisecond)
		}
	}
}

func (endpoint ReceiveEndpoint) TestInternal() {
	endpoint.MessageReceiver.Receive()
	endpoint.WaitForCompletion()
}
