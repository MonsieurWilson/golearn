// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 97.
//!+

// Charcount computes counts of Unicode characters.
package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"unicode"
)

func main() {
	counts := make(map[string]int) // counts of letters, numbers and marks etc.
	invalid := 0                   // count of invalid UTF-8 characters

	in := bufio.NewReader(os.Stdin)
	for {
		r, n, err := in.ReadRune() // returns rune, nbytes, error
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("charcount: %v\n", err)
		}
		if r == unicode.ReplacementChar && n == 1 {
			invalid++
			continue
		}

		switch {
		case unicode.IsControl(r):
			counts["Control"]++
		case unicode.IsDigit(r):
			counts["Digit"]++
		case unicode.IsLetter(r):
			counts["Letter"]++
		case unicode.IsNumber(r):
			counts["Number"]++
		case unicode.IsPunct(r):
			counts["Puncuation"]++
		case unicode.IsSymbol(r):
			counts["Symbol"]++
		default:
			counts["Others"]++
		}
	}
	fmt.Printf("rune\t\tcount\n")
	for c, n := range counts {
		fmt.Printf("%q\t%d\n", c, n)
	}
	if invalid > 0 {
		fmt.Printf("\n%d invalid UTF-8 characters\n", invalid)
	}
}

//!-
