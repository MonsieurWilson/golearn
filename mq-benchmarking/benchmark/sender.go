package benchmark

import (
	"log"
	"sync/atomic"
	"time"
)

type MessageSender interface {
	Send([]byte, bool, func(string, error))
	Teardown()
}

type SendEndpoint struct {
	MessageSender    MessageSender
	NumberOfMessages int
	SizeOfMessage    int
	Async            bool
	Random           bool
}

func NewSendEndpoint(tester Tester) *SendEndpoint {
	return &SendEndpoint{
		MessageSender:    tester.Sender,
		NumberOfMessages: tester.MessageCount,
		SizeOfMessage:    tester.MessageSize,
		Async:            tester.Async,
		Random:           tester.Random,
	}
}

func (endpoint SendEndpoint) TestInternal() {
	var nAck int64
	message := make([]byte, endpoint.SizeOfMessage)
	done := make(chan bool)
	async := endpoint.Async
	random := endpoint.Random

	if random {
		RandBytes(message)
	}
	cb := func(string, error) {
		nAck = atomic.AddInt64(&nAck, 1)
		if nAck == int64(endpoint.NumberOfMessages) {
			done <- true
		}
	}
	start := time.Now().UnixNano()
	for i := 0; i < endpoint.NumberOfMessages; i++ {
		endpoint.MessageSender.Send(message, async, cb)
	}
	if async {
		<-done
	}
	stop := time.Now().UnixNano()
	elapsed := float32(stop-start) / 1e9
	log.Printf("Sent %d messages (size %d) in %.2f s",
		endpoint.NumberOfMessages,
		endpoint.SizeOfMessage,
		elapsed)
	log.Printf("Rate %.2f Mbps, %.2f msg/s",
		float32(8*endpoint.SizeOfMessage*endpoint.NumberOfMessages)/elapsed/1e6,
		float32(endpoint.NumberOfMessages)/elapsed)
}
