// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 73.

// Comma prints its argument numbers with a comma at each power of 1000.
//
// Example:
// 	$ go build gopl.io/ch3/comma
//	$ ./comma 1 12 123 1234 1234567890
// 	1
// 	12
// 	123
// 	1,234
// 	1,234,567,890
//
package main

import (
	"bytes"
	"fmt"
	"os"
	"strings"
)

func main() {
	for i := 1; i < len(os.Args); i++ {
		s := strings.Replace(os.Args[i], " ", "", -1)
		fmt.Printf("Comma representation: %s\n", comma(s))
	}
}

//!+
// comma inserts commas in a non-negative decimal integer string.
func comma(s string) string {
	var buf bytes.Buffer

	// s contains sign +/-.
	if strings.HasPrefix(s, "+") || strings.HasPrefix(s, "-") {
		buf.WriteString(s[:0])
	}

	// s represents a float number.
	if i := strings.Index(s, "."); i != -1 {
		buf.WriteString(comma(s[:i]))
		buf.WriteString(".")
		buf.WriteString(comma(s[i+1:]))
		return buf.String()
		// return strings.Join([]string{comma(s[:i]), comma(s[i+1:])}, ".")
	}

	// s represents a integer.
	lens := len(s)
	firstCommaPos := lens % 3
	buf.WriteString(s[:firstCommaPos])
	for i := firstCommaPos; i <= lens-3; i += 3 {
		if i != 0 {
			buf.WriteString(",")
		}
		buf.WriteString(s[i : i+3])
	}
	return buf.String()
}

//!-
