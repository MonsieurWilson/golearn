package main

import (
	"os"

	"github.com/nats-io/mq-benchmarking/benchmark"
	"github.com/nats-io/mq-benchmarking/benchmark/mq"
)

func main() {
	opt := benchmark.ParseArgs(os.Args)
	tester := newTester(opt)

	var msgQ benchmark.MessageQueue
	switch tester.Name {
	case "rabbitmq":
		msgQ = mq.NewRabbitmq(opt)
	case "nsq":
		msgQ = mq.NewNsq(opt)
	case "nats-streaming":
		msgQ = mq.NewStan(opt)
	case "nats":
		msgQ = mq.NewGnatsd(opt)
	default:
		benchmark.Usage()
	}

	tester.Test(msgQ)
}

func newTester(opt *benchmark.Options) *benchmark.Tester {
	return &benchmark.Tester{
		Name:         opt.Mq,
		MessageSize:  opt.MessageSize,
		MessageCount: opt.MessageCount,
		Async:        opt.Async,
		Random:       opt.Random,
		Direction:    opt.Direction,
		Core:         opt.Core,
	}
}
