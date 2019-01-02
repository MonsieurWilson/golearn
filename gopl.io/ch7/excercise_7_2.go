// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 173.

// ByteCountingWriter demonstrates an implementation of io.Writer that
// wraps a given writer and a counter implies how many bytes have written.
package main

import (
	"fmt"
	"io"
	"os"
)

func CountingWriter(w io.Writer) (io.Writer, *int64) {
	bcw := ByteCountingWriter{w, new(int64)}
	return bcw, bcw.counter
}

//!+bytecoutingwriter

type ByteCountingWriter struct {
	w       io.Writer
	counter *int64
}

func (bcw ByteCountingWriter) Write(d []byte) (int, error) {
	*bcw.counter += int64(len(d))
	return bcw.w.Write(d)
}

//!-bytecoutingwriter

func main() {
	//!+main
	bcw, counter := CountingWriter(os.Stdout)
	bcw.Write([]byte("world\n"))
	fmt.Println("Total have written:", *counter)

	bcw.Write([]byte("hello\n"))
	fmt.Println("Total have written:", *counter)
	//!-main
}
