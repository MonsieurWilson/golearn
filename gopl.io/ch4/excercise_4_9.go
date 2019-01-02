package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	input := bufio.NewScanner(os.Stdin)
	counts := make(map[string]int)
	input.Split(bufio.ScanWords)

	for input.Scan() {
		counts[input.Text()]++
	}

	if err := input.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "Reading input error: ", err)
	} else {
		fmt.Fprintln(os.Stdout, "-------------")
		fmt.Fprintln(os.Stdout, "Word\t\tCount")
		for key, val := range counts {
			fmt.Fprintf(os.Stdout, "%s\t\t%d\n", key, val)
		}
	}
}
