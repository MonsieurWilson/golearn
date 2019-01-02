package mq

import (
	"log"

	"github.com/nats-io/go-nats"
	"github.com/nats-io/mq-benchmarking/benchmark"
)

type Gnatsd struct {
	handler benchmark.MessageHandler
	conn    *nats.Conn
	topic   string
}

func NewGnatsd(opt *benchmark.Options) *Gnatsd {
	url := nats.DefaultURL
	if opt.Server != "" {
		url = opt.Server
	}

	var conn *nats.Conn
	var err error
	switch {
	case opt.CaFile != "":
		conn, err = nats.Connect(url, nats.RootCAs(opt.CaFile))
	case opt.CertFile != "" && opt.KeyFile != "":
		conn, err = nats.Connect(url, nats.ClientCert(opt.CertFile, opt.KeyFile))
	case opt.CertFile == "" && opt.KeyFile == "":
		conn, err = nats.Connect(url)
	default:
		log.Fatalf("Error: tls_cert and tls_key should be configured together.")
	}
	if err != nil {
		log.Fatalf("Can't connect NATS server: %v.", err)
	}

	// We want to be alerted if we get disconnected, this will
	// be due to Slow Consumer.
	conn.Opts.AllowReconnect = false

	// Report async errors.
	conn.Opts.AsyncErrorCB = func(nc *nats.Conn, sub *nats.Subscription, err error) {
		log.Fatalf("NATS: Received an async error! %v", err)
	}

	// Report a reconnect scenario.
	conn.Opts.ReconnectedCB = func(nc *nats.Conn) {
		log.Fatalf("NATS: Reconnected NATS server!")
	}

	var handler benchmark.MessageHandler
	handler = &benchmark.ThroughputMessageHandler{
		NumberOfMessages: opt.MessageCount,
		SizeOfMessage:    opt.MessageSize,
	}

	return &Gnatsd{
		handler: handler,
		topic:   opt.Topic,
		conn:    conn,
	}
}

func (g *Gnatsd) Teardown() {
	g.conn.Flush()
	g.conn.Close()
}

func (g *Gnatsd) Receive() {
	g.conn.Subscribe(g.topic, func(message *nats.Msg) {
		g.handler.ReceiveMessage(message.Data)
	})
}

func (g *Gnatsd) Send(message []byte, async bool, cb func(string, error)) {
	g.conn.Publish(g.topic, message)
}

func (g *Gnatsd) MessageHandler() *benchmark.MessageHandler {
	return &g.handler
}
