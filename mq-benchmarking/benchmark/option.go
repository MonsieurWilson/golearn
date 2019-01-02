package benchmark

import (
	"flag"
	"fmt"
	"log"
	"os"
)

type Options struct {
	Mq           string
	Async        bool
	Random       bool
	MessageCount int
	MessageSize  int
	Server       string
	Topic        string
	Direction    string
	CaFile       string
	CertFile     string
	KeyFile      string
	Core         string
}

var usageStr = `
Usage: mq-benchmarking <message queue> [options]

Message Queue Options:
    rabbitmq       - RabbitMQ
    nsq            - NSQ
    nats           - NATS
    nats-streaming - NATS Streaming

Benchmarking Options:
    -n,  --num_messages <int> Number of total messages (default: 1000000)
    -s,  --message_size <int> Size of message (default: 1000)
    -a,  --async              Async message publishing (default: false)
    -r,  --random             Send random bytes instead of empty bytes (default: false)
    -u,  --url <string>       Message queue server URL (default: varies from different message queues)
    -t,  --topic <string>     Topic of publish/subscribe (default: test)
    -d,  --direction <string> Direction of test, eg: send|receive (default: sender)

TLS Options:
    --tls_cert    <path>      Client certificate file
    --tls_key     <path>      Client key file
    --tls_cacert  <path>      Client CA certificate file

Other Options:
    -c, --core <string>       Set CPU affinity, eg: 1|2-3 (default: all available CPUs)
`

func Usage() {
	fmt.Printf("%s\n", usageStr)
	os.Exit(0)
}

func ParseArgs(args []string) *Options {
	if len(args) < 2 {
		Usage()
	}

	var mq = args[1]
	var messageCount, messageSize int
	var server, topic, direction string
	var async, random bool
	var certFile, keyFile, caFile string
	var core string

	fs := flag.NewFlagSet("benchmarking", flag.ExitOnError)
	fs.Usage = Usage
	fs.IntVar(&messageCount, "n", 1000000, "Number of total messages")
	fs.IntVar(&messageCount, "num_messages", 1000000, "Number of total messages")
	fs.IntVar(&messageSize, "s", 1000, "Size of message")
	fs.IntVar(&messageSize, "message_size", 1000, "Size of message")
	fs.BoolVar(&async, "a", false, "Async message publishing")
	fs.BoolVar(&async, "async", false, "Async message publishing")
	fs.BoolVar(&random, "r", false, "Send random bytes instead of empty bytes")
	fs.BoolVar(&random, "random", false, "Send random bytes instead of empty bytes")
	fs.StringVar(&server, "u", "", "Message queue server")
	fs.StringVar(&server, "url", "", "Message queue server")
	fs.StringVar(&topic, "t", "test", "Publish/subscribe topic")
	fs.StringVar(&topic, "topic", "test", "Publish/subscribe topic")
	fs.StringVar(&direction, "d", "send", "Direction of test")
	fs.StringVar(&direction, "direction", "send", "Direction of test")
	// TLS options
	fs.StringVar(&certFile, "tls_cert", "", "Client certificate file")
	fs.StringVar(&keyFile, "tls_key", "", "Client key file")
	fs.StringVar(&caFile, "tls_cacert", "", "Client CA certificate file")
	fs.StringVar(&core, "c", "", "Set CPU affinity")
	fs.StringVar(&core, "core", "", "Set CPU affinity")

	fs.Parse(args[2:])

	if async && mq != "nats-streaming" {
		log.Println("Only nats-streaming support async benchmark at present.")
		os.Exit(1)
	}

	return &Options{
		Mq:           mq,
		MessageCount: messageCount,
		MessageSize:  messageSize,
		Async:        async,
		Random:       random,
		Server:       server,
		Topic:        topic,
		Direction:    direction,
		CaFile:       caFile,
		CertFile:     certFile,
		KeyFile:      keyFile,
		Core:         core,
	}
}
