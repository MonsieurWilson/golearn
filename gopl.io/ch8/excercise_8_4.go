// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 224.

// Reverb2 is a TCP server that simulates an echo.
package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
	"sync"
	"time"
)

//!+
func closeConn(conn net.Conn) {
	if x, ok := conn.(*net.TCPConn); ok {
		x.CloseWrite()
	} else {
		x.Close()
	}
}

func echo(conn net.Conn, shout string, delay time.Duration, wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Fprintln(conn, "\t", strings.ToUpper(shout))
	time.Sleep(delay)
	fmt.Fprintln(conn, "\t", shout)
	time.Sleep(delay)
	fmt.Fprintln(conn, "\t", strings.ToLower(shout))
}

func handleConn(conn net.Conn) {
	defer closeConn(conn)
	var wg sync.WaitGroup
	input := bufio.NewScanner(conn)
	for input.Scan() {
		wg.Add(1)
		go echo(conn, input.Text(), 1*time.Second, &wg)
	}
	// NOTE: ignoring potential errors from input.Err()

	wg.Wait()
}

//!-

func main() {
	l, err := net.Listen("tcp", "0.0.0.0:8001")
	if err != nil {
		log.Fatal(err)
	}
	for {
		conn, err := l.Accept()
		if err != nil {
			log.Print(err) // e.g., connection aborted
			continue
		}
		go handleConn(conn)
	}
}
