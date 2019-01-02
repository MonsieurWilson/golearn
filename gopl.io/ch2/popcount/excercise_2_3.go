// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 45.

// (Package doc comment intentionally malformed to demonstrate golint.)
//!+
package popcount

// pc[i] is the population count of i.
var pc [256]byte

// PopCount returns the population count (number of set bits) of x.
func PopCount(x uint64) int {
	var count int
	for i := uint64(0); i < 8; i++ {
		count += int(pc[byte(x>>8*i)])
	}
	return count
}

func init() {
	for i := range pc {
		pc[i] = pc[i/2] + byte(i&1)
	}
}

//!-
