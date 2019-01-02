package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"runtime"
	"sync/atomic"
	"syscall"

	"github.com/containerd/cgroups"
	nats "github.com/nats-io/go-nats"
	"github.com/nats-io/go-nats-streaming"
	"github.com/nats-io/nats-streaming-server/server"
	specs "github.com/opencontainers/runtime-spec/specs-go"
)

const (
	DefaultServerAddr     = "nats://172.16.0.1:4222"
	DefaultUnixSockPath   = "/var/run/npai-agent.sock"
	DefaultMsgChannelSize = 1024
)

var usageStr = `
Usage: npai-agent [options] <message queue> <topic>

Message Queue Options:
    nats           - NATS
    nats-streaming - NATS streaming

Options:
    -s, --server  <url>             NATS server URL (default: nats://127.0.0.1:4222)
    -l, --listen  <url>             Agent listen URL (default: /var/run/npai-agent.sock)
    -c, --core    <string>          Set CPU affinity, eg: 1|2-3 (default: all available CPUs)

TLS Options:
    --tls_cert    <path>            Client certificate file
    --tls_key     <path>            Client key file
    --tls_cacert  <path>            Client CA certificate file
`

func usage() {
	log.Fatalf(usageStr)
}

type Msg struct {
	Data []byte
}

func main() {
	var serverAddr, unixSockPath string
	var core string
	var certFile, keyFile, caFile string
	flag.StringVar(&serverAddr, "s", DefaultServerAddr, "The NATS server URLs (separated by comma)")
	flag.StringVar(&serverAddr, "server", DefaultServerAddr, "The NATS server URLs (separated by comma)")
	flag.StringVar(&unixSockPath, "l", DefaultUnixSockPath, "Agent listen URL")
	flag.StringVar(&unixSockPath, "listen", DefaultUnixSockPath, "Agent listen URL")
	flag.StringVar(&core, "c", "", "Set CPU affinity")
	flag.StringVar(&core, "core", "", "Set CPU affinity")
	flag.StringVar(&certFile, "tls_cert", "", "Client certificate file")
	flag.StringVar(&keyFile, "tls_key", "", "Client key file")
	flag.StringVar(&caFile, "tls_cacert", "", "Client CA certificate file")

	log.SetFlags(0)
	flag.Usage = usage
	flag.Parse()

	args := flag.Args()
	if len(args) < 2 {
		log.Printf("Error: Message queue and topic mush be specified.")
		usage()
	}

	var err error
	mq, topic := args[0], args[1]

	// Configure Cgroups
	memNode := "0"
	if core == "" {
		core = fmt.Sprintf("0-%d", runtime.NumCPU()-1)
	}
	cgroupName := fmt.Sprintf("/%s-npai-agent", mq)

	cgCtrl, err := cgroups.New(cgroups.V1, cgroups.StaticPath(cgroupName),
		&specs.LinuxResources{
			CPU: &specs.LinuxCPU{
				Cpus: core,
				Mems: memNode,
			},
		})
	if err != nil {
		log.Fatalf("Error: failed to create cgroup: %v", err)
	}
	if err := cgCtrl.Add(cgroups.Process{Pid: os.Getpid()}); err != nil {
		log.Fatalf("Error: failed to add process to cgroup: %v", err)
	}

	// Connect NATS server
	var nc *nats.Conn
	switch {
	case caFile != "":
		nc, err = nats.Connect(serverAddr, nats.RootCAs(caFile))
	case certFile != "" && keyFile != "":
		nc, err = nats.Connect(serverAddr, nats.ClientCert(certFile, keyFile))
	case certFile == "" && keyFile == "":
		nc, err = nats.Connect(serverAddr)
	default:
		log.Fatalf("Error: tls_cert and tls_key should be configured together.")
		usage()
	}
	if err != nil {
		log.Fatalf("Can't connect: %v\n", err)
	}
	defer nc.Close()

	// Setup message handlers
	var deliverMsgs func()
	msgChan := make(chan *Msg, DefaultMsgChannelSize)
	switch mq {
	case "nats":
		deliverMsgs = func() {
			for msg := range msgChan {
				nc.Publish(topic, msg.Data)
			}
		}
	case "nats-streaming":
		// connect Nats-streaming server first
		sc, err := stan.Connect(server.DefaultClusterID, "npai-agent",
			stan.NatsConn(nc))
		if err != nil {
			log.Fatalf("Can't connect server: %v\n", err)
		}

		// using async publish method
		deliverMsgs = func() {
			for msg := range msgChan {
				sc.PublishAsync(topic, msg.Data, nil)
			}
		}
	default:
		log.Fatalf("Unsupported message queue")
	}

	// Sit on a unix socket
	listener, err := net.Listen("unix", unixSockPath)
	log.Printf("Listening on [%v]", unixSockPath)

	// Spawn a goroutine to handle signals
	var messageCount uint64
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, os.Kill, syscall.SIGTERM)
	go func() {
		for range sigChan {
			log.Printf("Total received messages: %v", messageCount)

			// clean unix socket
			syscall.Unlink(unixSockPath)

			// clean Cgroups
			rootCgCtrl, err := cgroups.Load(cgroups.V1, cgroups.RootPath)
			if err != nil {
				log.Fatalf("Error: failed to load root cgroup: %v", err)
			}
			cgCtrl.MoveTo(rootCgCtrl)
			cgCtrl.Delete()

			os.Exit(0)
		}
	}()

	// Spawn a goroutine to deliver messages to broker
	go deliverMsgs()

	// Loop
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatalf("Error: failed to accept connection: %v", err)
		}

		go func(conn net.Conn) {
			scanner := bufio.NewScanner(conn)
			scanner.Split(bufio.ScanLines)
			for scanner.Scan() {
				msgChan <- &Msg{scanner.Bytes()}
				atomic.AddUint64(&messageCount, 1)
			}
		}(conn)
	}
}
