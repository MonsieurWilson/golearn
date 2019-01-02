#!/usr/bin/env python
# -*- coding: utf-8 -*-
# __author__ = "Wilson Lan"

import socket
import time
import sys
import os
import argparse
import errno
import multiprocessing

# NOTE:
#   In some cases, cgroups.py doesn't locate in current directory.
#   Callers should be aware of this situation and expand the syspath, 
#   e.g. sys.path.append("your_own_python_lib_path")
sys.path.append("/home/lanxinyu/nats/pylibs")
import cgroups

DEFAULT_NUM_MSGS      = 1000000
DEFAULT_MSG_SIZE      = 256
DEFAULT_NUM_PROCESSES = 2


def show_usage():
    message = """
Usage: npaiagent_test [options]

options:
    -n, --count     COUNT              Messages to send (default: 1000000)
    -s, --size      SIZE               Message size (default: 256)
    -u, --url       URL                NPAI agent address (default: /var/run/npai-agent.sock)
    -c, --core      CORE               Set CPU affinity (default: all available CPUs)
    -p, --processes NUM                Number of subprocess to send messages (default: 2)
    """
    print(message)

def show_usage_and_die():
    show_usage()
    sys.exit(1)

def configure_cgroups(core):
    cg = []
    cg.append(cgroups.must_make_cgroup("cpuset", "npaiagent_test"))
    cg.append(cgroups.must_make_cgroup("cpu", "npaiagent_test"))
    
    # set cpuset attributes
    if core == "":
        cgroups.set_cgroups_cpuset_attr(*cg)
    else:
        cgroups.set_cgroups_cpuset_attr(*cg, core_id=core)

    # set cgroup process
    cgroups.set_cgroups_procs(os.getpid(), *cg)

    return cg

def send_message(url, count, size):
    # Initialize a unix socket
    sock = socket.socket(socket.AF_UNIX, socket.SOCK_STREAM)

    # Connect to npai-agent 
    returncode = sock.connect_ex(url)
    if returncode != 0:
        print "Error: failed to connect {url}: {errmsg}".format(
                url=url, errmsg=os.strerror(returncode))
        return

    # Formulate packet payload
    data = []
    for i in range(0, size):
        data.append(b'W')
    data.append(b'\n')

    payload = b''.join(data)

    # Deliver messages
    for i in xrange(count):
        # NOTE Using `sendall` will lead to poor performance.
        # Perhaps, the calling of `send` over unix socket is okay.
        # However, we should write like this in those serious occasions:
        ## sent = 0
        ## while sent < size+1:
        ##     nbytes = sock.send(payload[sent:])
        ##     sent += nbytes
        sock.send(payload)

    sock.close()

def main():
    parser = argparse.ArgumentParser()
    parser.add_argument("-n", "--count", default=DEFAULT_NUM_MSGS, type=int)
    parser.add_argument("-s", "--size", default=DEFAULT_MSG_SIZE, type=int)
    parser.add_argument("-u", "--url", default="/var/run/npai-agent.sock")
    parser.add_argument("-c", "--core", default="")
    parser.add_argument("-p", "--processes", default=DEFAULT_NUM_PROCESSES,
            type=int)
    args = parser.parse_args()

    # Configure Cgroups
    cg = configure_cgroups(args.core)

    # Start the benchmark
    start = time.time()
    p = multiprocessing.Pool()
    for i in xrange(args.processes):
        p.apply_async(send_message, args=(args.url, 
            args.count/args.processes, args.size))
    p.close()
    p.join()
    elapsed = time.time() - start

    # Print throughput
    throughput = args.count / elapsed
    rate = args.size * args.count * 8 / elapsed / 1e6
    print "Send %d messages (size %d) in %.2f seconds" % \
        (args.count, args.size, elapsed)
    print "Rate %.2f msg/s, %.2f Mbps" % (throughput, rate)

    # Clean Cgroups
    cgroups.clean_cgroups(*cg)


if __name__ == "__main__":
    main()

