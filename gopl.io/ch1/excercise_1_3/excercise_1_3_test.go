package main

import "testing"

var strlist []string

func BenchmarkStrjoin(b *testing.B) {
	for i := 0; i < b.N; i++ {
		strjoin(strlist, " ")
	}
}

func BenchmarkPlusjoin(b *testing.B) {
	for i := 0; i < b.N; i++ {
		plusjoin(strlist, " ")
	}
}

func init() {
	strlist = make([]string, 1000)
	for i, _ := range strlist {
		strlist[i] = "W"
	}
}
