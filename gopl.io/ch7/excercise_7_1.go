// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 173.

// Bytecounter demonstrates an implementation of io.Writer that counts bytes.
package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"strings"
)

//!+bytecounter

type WordCounter int64

func (c *WordCounter) Write(p []byte) (int64, error) {
	reader := strings.NewReader(string(p))
	scanner := bufio.NewScanner(reader)
	scanner.Split(bufio.ScanWords)

	for scanner.Scan() {
		*c += WordCounter(1)
	}

	return int64(*c), nil
}

func (c *WordCounter) String() string {
	return fmt.Sprint(int64(*c))
}

//!-bytecounter

func main() {
	//!+main
	var c WordCounter
	var input bytes.Buffer
	input.ReadFrom(os.Stdin)
	c.Write(input.Bytes())
	fmt.Println(c) // "5", = len("hello")

	//!-main
}
