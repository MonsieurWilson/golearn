package main

import "fmt"

const (
	B  = 1
	KB = 1000 * B
	MB = 1000 * KB
	GB = 1000 * MB
	TB = 1000 * GB
	PB = 1000 * TB
	EB = 1000 * PB
	ZB = 1000 * EB
	YB = 1000 * ZB
)

func main() {
	fmt.Printf("%v\n%v\n%v\n%v\n%v\n%v\n%v\n",
		B, KB, MB, GB, TB, PB, EB)
	fmt.Printf("%v\n", YB/ZB)
}
