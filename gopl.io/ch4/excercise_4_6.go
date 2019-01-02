package main

import (
	"bytes"
	"fmt"
	"unicode"
)

func main() {
	strings := [][]byte{
		[]byte("  hello   world    "),
		[]byte("   "),
		[]byte("hello  "),
		[]byte(" hel  lo  "),
	}
	for _, s := range strings {
		fmt.Println("=========================================")
		fmt.Printf("%23v: $%v$\n", "original", string(s))
		var buf1, buf2 bytes.Buffer
		buf1.Write(s)
		buf2.Write(s)
		fmt.Printf("%23v: $%v$\n", "delete duplicate spaces", string(delDupSpaces(buf1.Bytes())))
		fmt.Printf("%23v: $%v$\n", "delete all spaces", string(delSpaces(buf2.Bytes())))
	}
}

// delete spaces from byte array.
func delSpaces(s []byte) []byte {
	lens := len(s)
	var head int
	for cur := head; cur < lens; cur, head = cur+1, head+1 {
		// ignore space byte.
		for cur < lens && unicode.IsSpace(rune(s[cur])) {
			cur++
		}
		if cur >= lens {
			break
		}
		s[head] = s[cur]
	}
	return s[:head]
}

// delete duplicate spaces from byte array.
func delDupSpaces(s []byte) []byte {
	lens := len(s)
	if lens < 2 {
		return s
	}

	var head int
	for cur := head; cur < lens; cur, head = cur+1, head+1 {
		// ignore space byte.
		if unicode.IsSpace(rune(s[cur])) {
			for cur < lens-1 && s[cur] == s[cur+1] {
				cur++
			}
		}
		s[head] = s[cur]
	}
	return s[:head]
}
