package main

import "fmt"

func main() {
	strings := [][]string{
		[]string{"I'm", "I'm", "dup", "dup", "helo", "hello", "centos", "centos"},
		[]string{"I'm"},
		[]string{"I'm", "I'm"},
		[]string{"I'm", "I'm", "I'm", "Good"},
	}
	for _, s := range strings {
		fmt.Printf("%v\n", delDuplicate(s))
	}
}

// delete duplicate elements from string array.
func delDuplicate(s []string) []string {
	lens := len(s)
	if lens < 2 {
		return s
	}

	var head int
	for cur := head; cur < lens; cur, head = cur+1, head+1 {
		for cur < lens-1 && s[cur] == s[cur+1] {
			cur++
		}
		s[head] = s[cur]
	}
	return s[:head]
}
