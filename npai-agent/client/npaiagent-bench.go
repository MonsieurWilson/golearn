package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"runtime"
	"time"

	"github.com/containerd/cgroups"
	"github.com/nats-io/mq-benchmarking/benchmark"
	specs "github.com/opencontainers/runtime-spec/specs-go"
)

const (
	DefaultUnixSockPath = "/var/run/npai-agent.sock"
	DefaultMessageCount = 1000000
	DefaultMessageSize  = 1000
)

var usageStr = `
Usage npaiagent_test [options]

Options:
    -u, --url  <url>           Npai agent URL (default: /var/run/npai-agent.sock)
    -c, --core <string>        Set CPU afinity, eg: 1|2-3 (default: all available CPUs)
    -n, --num_messages <int>   Number of total messages (default: 1000000)
    -s, --message_size <int>   Size of message (default: 1000)
    -r, --random               Send random bytes instead of empty bytes (default: false)
`

func usage() {
	log.Fatalf(usageStr)
}

func main() {
	var url string
	var core string
	var messageCount int
	var messageSize int
	var random bool
	flag.StringVar(&url, "u", DefaultUnixSockPath, "Npai agent URL")
	flag.StringVar(&url, "url", DefaultUnixSockPath, "Npai agent URL")
	flag.StringVar(&core, "c", "", "Set CPU afinity")
	flag.StringVar(&core, "core", "", "Set CPU afinity")
	flag.IntVar(&messageCount, "n", DefaultMessageCount, "Number of total messages")
	flag.IntVar(&messageCount, "num_messages", DefaultMessageCount, "Number of total messages")
	flag.IntVar(&messageSize, "s", DefaultMessageSize, "Size of message")
	flag.IntVar(&messageSize, "message_size", DefaultMessageSize, "Size of message")
	flag.BoolVar(&random, "r", false, "Send random bytes instead of empty bytes")
	flag.BoolVar(&random, "random", false, "Send random bytes instead of empty bytes")
	log.SetFlags(0)
	flag.Usage = usage
	flag.Parse()

	// Configure Cgroups
	memNode := "0"
	if core == "" {
		core = fmt.Sprintf("0-%d", runtime.NumCPU()-1)
	}
	cgroupName := "/npaiagent_test"
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

	// Connect to npai-agent
	conn, err := net.Dial("unix", url)
	if err != nil {
		log.Fatalf("Error: Can't connect to npai-agent: %v", err)
	}

	// Formulate payload
	message := make([]byte, messageSize)
	if random {
		benchmark.RandBytes(message)
	}
	message = append(message, '\n')

	// Start benchmarking
	start := time.Now().UnixNano()
	for i := 0; i < messageCount; i++ {
		conn.Write(message)
	}
	stop := time.Now().UnixNano()
	elapsed := float32(stop-start) / 1e9
	log.Printf("Sent %d messages (size %d) in %.2f s",
		messageCount,
		messageSize,
		elapsed)
	log.Printf("Rate %.2f Mbps, %.2f msg/s",
		float32(8*messageSize*messageCount)/elapsed/1e6,
		float32(messageCount)/elapsed)

	// Close connection
	conn.Close()

	// Clean Cgroups
	rootCgCtrl, err := cgroups.Load(cgroups.V1, cgroups.RootPath)
	if err != nil {
		log.Fatalf("Error: failed to load root cgroup: %v", err)
	}
	cgCtrl.MoveTo(rootCgCtrl)
	cgCtrl.Delete()
}
