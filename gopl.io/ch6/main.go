package main

import (
	"fmt"

	"nsfocus/Wilson/gopl.io/ch6/intset"
	"nsfocus/Wilson/gopl.io/ch6/intset2"
)

func main() {
	var set intset.IntSet

	fmt.Println("Arch:", intset2.OpMode)

	fmt.Println("Original set:")
	set.Add(1)
	set.AddAll(2, 3)
	fmt.Printf("%v\n===============\n", &set)

	fmt.Println("Create a copy:")
	copyset := *set.Copy()
	fmt.Printf("%v\n===============\n", &copyset)

	fmt.Println("Before remove:")
	fmt.Printf("%v\n===============\n", &set)

	// remove elements
	set.Remove(1)
	set.Remove(2)

	fmt.Println("After remove 1 and 2:")
	fmt.Printf("%v\n===============\n", &set)

	// clear set
	set.Clear()

	fmt.Println("After clear:")
	fmt.Printf("%v\n===============\n", &set)

	fmt.Println("The copy set:")
	fmt.Printf("%v\n===============\n", &copyset)
}
