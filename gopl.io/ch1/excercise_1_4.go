// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 10.
//!+

// Dup2 prints the count and text of lines that appear more than once
// in the input.  It reads from stdin or from a list of named files.
package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	dupLineFiles := make(map[string]struct{})
	files := os.Args[1:]
	if len(files) == 0 {
		countLines(os.Stdin, dupLineFiles)
	} else {
		for _, arg := range files {
			f, err := os.Open(arg)
			if err != nil {
				fmt.Fprintf(os.Stderr, "dup2: %v\n", err)
				continue
			}
			countLines(f, dupLineFiles)
			f.Close()
		}
	}

	fmt.Println("=====\nFiles contain duplicate lines:")
	for filename, _ := range dupLineFiles {
		fmt.Printf("%v\n", filename)
	}
}

func countLines(f *os.File, dupLineFiles map[string]struct{}) {
	counts := make(map[string]int)
	input := bufio.NewScanner(f)
	for input.Scan() {
		text := input.Text()
		counts[text]++
		if counts[text] > 1 {
			dupLineFiles[f.Name()] = struct{}{}
		}
	}
	// NOTE: ignoring potential errors from input.Err()
}

//!-
