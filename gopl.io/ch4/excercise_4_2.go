// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 83.

// The sha256 command computes the SHA256 hash (an array) of a string.
package main

import (
	"bufio"
	"bytes"
	"crypto/sha256"
	"crypto/sha512"
	"flag"
	"log"
	"os"
)

//!+

var usageStr = `
Usage: excercise_4_2 [options]

Options:
    --ctype <type>                Crypto type, sha256(default), sha384, sha512
`

func usage() {
	log.Fatalf(usageStr)
}

func main() {
	// Clear log output format
	log.SetFlags(0)

	var cryptoType string
	flag.StringVar(&cryptoType, "ctype", "sha256", "crypto type")
	flag.Usage = usage
	flag.Parse()

	inputScanner := bufio.NewScanner(os.Stdin)
	log.Printf("Input strings:")
	for inputScanner.Scan() {
		var outputBuffer bytes.Buffer
		input := inputScanner.Bytes()
		switch cryptoType {
		case "sha256":
			shacode := sha256.Sum256(input)
			outputBuffer.Write(shacode[:])
		case "sha384":
			shacode := sha512.Sum384(input)
			outputBuffer.Write(shacode[:])
		case "sha512":
			shacode := sha512.Sum512(input)
			outputBuffer.Write(shacode[:])
		default:
			log.Fatalf("Unsupported crypto type!")
		}
		log.Printf("%v: %x\nInput strings:\n",
			cryptoType, outputBuffer.Bytes())
	}
}

//!-
