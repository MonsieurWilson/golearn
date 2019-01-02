// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 142.

// The sum program demonstrates a variadic function.
package main

import "fmt"

//!+
func max(vals ...int) (int, error) {
	if len(vals) == 0 {
		return 0, fmt.Errorf("function max's arguments can't be none")
	}
	maximum := vals[0]
	for _, val := range vals {
		if val > maximum {
			maximum = val
		}
	}
	return maximum, nil
}

func min(vals ...int) (int, error) {
	if len(vals) == 0 {
		return 0, fmt.Errorf("function min's arguments can't be none")
	}
	minimum := vals[0]
	for _, val := range vals {
		if val < minimum {
			minimum = val
		}
	}
	return minimum, nil
}

//!-

func main() {
	//!+main
	fmt.Println(max())           //  "0"
	fmt.Println(max(3))          //  "3"
	fmt.Println(max(1, 2, 3, 4)) //  "4"
	//!-main

	//!+slice
	values := []int{1, 2, 3, 4, 20}
	fmt.Println(max(values...)) // "20"
	//!-slice
}
