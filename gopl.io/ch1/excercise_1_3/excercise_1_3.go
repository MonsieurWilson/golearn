package main

import (
	"log"
	"strings"
	"time"
)

func strjoin(strlist []string, sep string) string {
	return strings.Join(strlist, sep)
}

func plusjoin(strlist []string, sep string) string {
	var res, s string
	for _, str := range strlist {
		res += s + str
		s = sep
	}
	return res
}

func main() {
	strlist := make([]string, 100000)
	for i, _ := range strlist {
		strlist[i] = "W"
	}

	var start, stop int64
	start = time.Now().UnixNano()
	strjoin(strlist, " ")
	stop = time.Now().UnixNano()
	log.Printf("Time elapsed (strings.Join version): %.4f s",
		float32(stop-start)/1e9)

	start = time.Now().UnixNano()
	plusjoin(strlist, " ")
	stop = time.Now().UnixNano()
	log.Printf("Time elapsed (\"+\" version): %.4f s",
		float32(stop-start)/1e9)
}
