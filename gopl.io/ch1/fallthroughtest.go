package main

import (
	"bufio"
	"log"
	"os"
)

func main() {
	log.SetFlags(0)

	scanner := bufio.NewScanner(os.Stdin)
	scanner.Split(bufio.ScanWords)
	for scanner.Scan() {
		switchtest(scanner.Text())
	}
}

func switchtest(num string) {
	switch num {
	case "0":
		log.Printf("drop in case 0!")
		fallthrough
	case "1":
		log.Printf("drop in case 1!")
		fallthrough
	default:
		log.Printf("drop in default!")
	}
}
