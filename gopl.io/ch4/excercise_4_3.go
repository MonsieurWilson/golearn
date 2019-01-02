// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 86.

// Rev reverses a slice.
package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	//!+array
	a := [6]int{0, 1, 2, 3, 4, 5}
	reverse(&a)
	fmt.Println(a) // "[5 4 3 2 1 0]"
	//!-array

	// Interactive test of reverse.
	fmt.Println("Input 6 numbers:")
	input := bufio.NewScanner(os.Stdin)
outer:
	for input.Scan() {
		idx := 0
		ints := [6]int{}
		for _, s := range strings.Fields(input.Text()) {
			x, err := strconv.ParseInt(s, 10, 64)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				continue outer
			}
			ints[idx] = int(x)
			idx++
		}
		reverse(&ints)
		fmt.Printf("After reverse: %v\n", ints)
	}
	// NOTE: ignoring potential errors from input.Err()
}

//!+rev
// reverse reverses an array of ints in place.
func reverse(a *[6]int) {
	for i, j := 0, 5; i < j; i, j = i+1, j-1 {
		a[i], a[j] = a[j], a[i]
	}
}

//!-rev
