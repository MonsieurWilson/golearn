// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 136.

// The toposort program prints the nodes of a DAG in topological order.
package main

import "fmt"

//!+table
// prereqs maps computer science courses to their prerequisites.
var prereqs = map[string]map[string]struct{}{
	"algorithms": {
		"data structures": struct{}{},
	},
	"calculus": {
		"linear algebra": struct{}{},
	},
	"compilers": {
		"data structures":       struct{}{},
		"formal languages":      struct{}{},
		"computer organization": struct{}{},
	},
	"data structures": {
		"discrete math": struct{}{},
	},
	"databases": {
		"data structures": struct{}{},
	},
	"discrete math": {
		"intro to programming": struct{}{},
	},
	"formal languages": {
		"discrete math": struct{}{},
	},
	"networks": {
		"operating systems": struct{}{},
	},
	"operating systems": {
		"data structures":       struct{}{},
		"computer organization": struct{}{},
	},
	"programming languages": {
		"data structures":       struct{}{},
		"computer organization": struct{}{},
	},
}

//!-table

//!+main
func main() {
	for i, course := range topoSort(prereqs) {
		fmt.Printf("%d:\t%s\n", i+1, course)
	}
}

func topoSort(m map[string]map[string]struct{}) []string {
	var order []string
	seen := make(map[string]bool)
	var visitAll func(items map[string]struct{})

	visitAll = func(items map[string]struct{}) {
		for item := range items {
			if !seen[item] {
				seen[item] = true
				visitAll(m[item])
				order = append(order, item)
			}
		}
	}

	items := make(map[string]struct{})
	for key := range m {
		items[key] = struct{}{}
	}

	visitAll(items)
	return order
}

//!-main
