package main

import (
	"fmt"
	"strings"
)

func main() {
	s := "hello world hello world"
	fmt.Println(expand(s, replace, "hello"))
}

func expand(s string, f func(string) string, p string) string {
	return strings.Replace(s, p, f(p), -1)
}

func replace(s string) string {
	return fmt.Sprintf("$$%s$$", s)
}
