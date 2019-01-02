package mq

import (
	"log"

	"github.com/nats-io/go-nats"
	"github.com/nats-io/go-nats-streaming"
	"github.com/nats-io/mq-benchmarking/benchmark"
	"github.com/nats-io/nats-streaming-server/server"
	"github.com/nats-io/nuid"
)

type Stan struct {
	handler benchmark.MessageHandler
	topic   string
	sc      stan.Conn
}

func NewStan(opt *benchmark.Options) *Stan {
	clientID := nuid.Next()
	url := stan.DefaultNatsURL
	if opt.Server != "" {
		url = opt.Server
	}

	var nc *nats.Conn
	var err error
	switch {
	case opt.CaFile != "":
		nc, err = nats.Connect(url, nats.RootCAs(opt.CaFile))
	case opt.CertFile != "" && opt.KeyFile != "":
		nc, err = nats.Connect(url, nats.ClientCert(opt.CertFile, opt.KeyFile))
	case opt.CertFile == "" && opt.KeyFile == "":
		nc, err = nats.Connect(url)
	default:
		log.Fatalf("Error: tls_cert and tls_key should be configured together.")
	}
	if err != nil {
		log.Fatalf("Can't connect NATS server: %v.", err)
	}

	sc, err := stan.Connect(server.DefaultClusterID, clientID, stan.NatsConn(nc))
	if err != nil {
		log.Fatalf("Can't connect Nats-streaming-server: %v.", err)
	}

	// We want to be alerted if we get disconnected, this will
	// be due to Slow Consumer.
	nc.Opts.AllowReconnect = true

	// Report async errors.
	nc.Opts.AsyncErrorCB = func(nc *nats.Conn, sub *nats.Subscription, err error) {
		log.Fatalf("NATS: Received an async error: %v!", err)
	}

	var handler benchmark.MessageHandler
	handler = &benchmark.ThroughputMessageHandler{
		NumberOfMessages: opt.MessageCount,
		SizeOfMessage:    opt.MessageSize,
	}

	return &Stan{
		handler: handler,
		topic:   opt.Topic,
		sc:      sc,
	}
}

func (s *Stan) Teardown() {
	s.sc.NatsConn().Close()
	s.sc.Close()
}

func (s *Stan) Receive() {
	startOpt := stan.DeliverAllAvailable()
	_, err := s.sc.Subscribe(s.topic, func(message *stan.Msg) {
		s.handler.ReceiveMessage(message.Data)
	}, startOpt)
	if err != nil {
		s.sc.Close()
		log.Fatalf("Failed to subscribe, %s", err)
	}
}

func (s *Stan) Send(message []byte, async bool, cb func(string, error)) {
	if async {
		s.sc.PublishAsync(s.topic, message, cb)
	} else {
		s.sc.Publish(s.topic, message)
	}
}

func (s *Stan) MessageHandler() *benchmark.MessageHandler {
	return &s.handler
}
