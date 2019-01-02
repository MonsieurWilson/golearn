package benchmark

import (
	"fmt"
	"log"
	"os"
	"runtime"

	"github.com/containerd/cgroups"
	specs "github.com/opencontainers/runtime-spec/specs-go"
)

type Tester struct {
	Name         string
	MessageSize  int
	MessageCount int
	Async        bool
	Random       bool
	Direction    string
	Core         string
	Sender       MessageSender
	Receiver     MessageReceiver
}

func (tester Tester) Test(msgQ MessageQueue) {
	// Configure Cgroups
	memNode := "0"
	if tester.Core == "" {
		tester.Core = fmt.Sprintf("0-%d", runtime.NumCPU()-1)
	}
	cgroupName := fmt.Sprintf("/%s-benchmark", tester.Name)

	cgCtrl, err := cgroups.New(cgroups.V1, cgroups.StaticPath(cgroupName),
		&specs.LinuxResources{
			CPU: &specs.LinuxCPU{
				Cpus: tester.Core,
				Mems: memNode,
			},
		})
	if err != nil {
		log.Fatalf("Error: failed to create cgroup: %v", err)
	}
	if err := cgCtrl.Add(cgroups.Process{Pid: os.Getpid()}); err != nil {
		log.Fatalf("Error: failed to add process to cgroup: %v", err)
	}

	// Start testing
	log.Printf("Begin %s test [PID %d]", tester.Name, os.Getpid())
	switch tester.Direction {
	case "send":
		log.Printf("========  Send test  ========")
		tester.Sender = msgQ
		sender := NewSendEndpoint(tester)
		sender.TestInternal()
		sender.MessageSender.Teardown()
	case "receive":
		log.Printf("======== Receive test ========")
		tester.Receiver = msgQ
		receiver := NewReceiveEndpoint(tester)
		receiver.TestInternal()
		receiver.MessageReceiver.Teardown()
	default:
		log.Fatalf("Unsupported direction.")
	}
	log.Printf("End %s test", tester.Name)

	// Clean Cgroups
	rootCgCtrl, err := cgroups.Load(cgroups.V1, cgroups.RootPath)
	if err != nil {
		log.Fatalf("Error: failed to load root cgroup: %v", err)
	}
	cgCtrl.MoveTo(rootCgCtrl)
	cgCtrl.Delete()
}
