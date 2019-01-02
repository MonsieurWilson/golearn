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
	a := []int{0, 1, 2, 3, 4, 5, 6}
	rotate(a, 2, true)
	fmt.Println(a) // "[2 3 4 5 6 0 1]"
	//!-array

	// Interactive test of reverse.
	input := bufio.NewScanner(os.Stdin)
outer:
	for input.Scan() {
		var ints []int
		for _, s := range strings.Fields(input.Text()) {
			x, err := strconv.ParseInt(s, 10, 64)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				continue outer
			}
			ints = append(ints, int(x))
		}
		rotate2(ints, 3, false)
		fmt.Printf("After reverse: %v\n", ints)
	}
	// NOTE: ignoring potential errors from input.Err()
}

//!+rotate
// rotate rotates a slice of ints.
func rotate(slice []int, step int, reverse bool) {
	if step >= len(slice) || step == 0 {
		return
	}

	rotateImpl(slice, step, reverse, 0)

	if len(slice)%2 == 0 && step%2 == 0 {
		rotateImpl(slice, step, reverse, 1)
	}
}

func rotateImpl(slice []int, step int, reverse bool, beg int) {
	bucket, length, sig := slice[0], len(slice), 1
	if reverse {
		sig = -1
	}
	for prev, next := beg, (length+sig*step)%length; ; prev, next = next, (next+length+sig*step)%length {
		if next == 0 {
			slice[prev] = bucket
			break
		} else {
			slice[prev] = slice[next]
		}
	}
}

// ??
func rotate2(slice []int, step int, reverse bool) {
	lens := len(slice)

	// no need to rotate.
	if step == 0 || lens == 0 {
		return
	}

	// step is too long, take the modulo value.
	if step >= lens {
		step %= lens
	}

	// rotate direction
	var sig int = 1
	if reverse {
		sig *= -1
	}

	for i := 0; i+step < lens; i++ {
		slice[i], slice[i+step] = slice[i+step], slice[i]
	}

	fmt.Println(lens - step)
	for i := lens - step; i < lens-1; i++ {
		slice[i], slice[(i+2*step)%lens] = slice[(i+2*step)%lens], slice[i]
	}
}

//!-rotate
