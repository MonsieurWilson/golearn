package main

import (
	"log"
	"os"
	"unicode"
	"unicode/utf8"
)

func main() {
	log.SetFlags(0)

	if len(os.Args) != 3 {
		log.Fatal("Invalid arguments")
	}
	log.Println("Result:",
		isAnagram([]rune(os.Args[1]), []rune(os.Args[2])))
}

func isAnagram(r1 []rune, r2 []rune) bool {
	if len(r1) != len(r2) {
		return false
	}

	dict := make(map[rune]int)
	const ch = utf8.RuneSelf
	for i := 0; i < len(r1); i++ {
		dict[unicode.ToLower(r1[i])-ch]++
		dict[unicode.ToLower(r2[i])-ch]--
	}
	for _, v := range dict {
		if v != 0 {
			return false
		}
	}
	return true
}
