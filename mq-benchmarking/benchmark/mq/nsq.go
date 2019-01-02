package mq

import (
	"github.com/bitly/go-nsq"
	"github.com/nats-io/mq-benchmarking/benchmark"
)

type Nsq struct {
	handler benchmark.MessageHandler
	pub     *nsq.Producer
	sub     *nsq.Consumer
	topic   string
	channel string
	server  string
}

func NewNsq(opt *benchmark.Options) *Nsq {
	topic := opt.Topic
	channel := opt.Topic
	url := "localhost:4150"
	if opt.Server != "" {
		url = opt.Server
	}

	pub, _ := nsq.NewProducer(url, nsq.NewConfig())
	sub, _ := nsq.NewConsumer(topic, channel, nsq.NewConfig())

	var handler benchmark.MessageHandler
	handler = &benchmark.ThroughputMessageHandler{
		NumberOfMessages: opt.MessageCount,
		SizeOfMessage:    opt.MessageSize,
	}

	return &Nsq{
		handler: handler,
		pub:     pub,
		sub:     sub,
		topic:   topic,
		channel: channel,
		server:  url,
	}
}

func (n *Nsq) Teardown() {
	n.sub.Stop()
	n.pub.Stop()
}

func (n *Nsq) Receive() {
	n.sub.AddHandler(nsq.HandlerFunc(func(message *nsq.Message) error {
		n.handler.ReceiveMessage(message.Body)
		return nil
	}))
	n.sub.ConnectToNSQD(n.server)
}

func (n *Nsq) Send(message []byte, async bool, cb func(string, error)) {
	n.pub.PublishAsync(n.topic, message, nil)
}

func (n *Nsq) MessageHandler() *benchmark.MessageHandler {
	return &n.handler
}
