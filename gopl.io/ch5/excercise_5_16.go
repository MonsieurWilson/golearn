package main

import (
	"fmt"
	"strings"
)

func StringJoin(sep string, strs ...string) string {
	var ret string
	for idx, val := range strs {
		ret += val
		if idx != len(strs)-1 {
			ret += sep
		}
	}
	return ret
}

func main() {
	strs := []string{"foo", "bar"}
	fmt.Println("StringJoin:", StringJoin(", ", strs...))
	fmt.Println("strings.Join:", strings.Join(strs, ", "))

	strs = []string{}
	fmt.Println("StringJoin:", StringJoin(", ", strs...))
	fmt.Println("strings.Join:", strings.Join(strs, ", "))

	strs = []string{"foo"}
	fmt.Println("StringJoin:", StringJoin(", ", strs...))
	fmt.Println("strings.Join:", strings.Join(strs, ", "))
}
