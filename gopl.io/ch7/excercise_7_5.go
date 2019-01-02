package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

func LimitReader(r io.Reader, n int64) io.Reader {
	return &LimitedReader{r, n}
}

//!+limitedreader

type LimitedReader struct {
	r io.Reader
	n int64
}

func (l *LimitedReader) Read(p []byte) (n int, err error) {
	if l.n <= 0 {
		return 0, io.EOF
	}
	if int64(len(p)) > l.n {
		p = p[0:l.n]
	}
	n, err = l.r.Read(p)
	l.n -= int64(n)
	return
}

//!-limitedreader

func main() {
	ltdreader := LimitReader(os.Stdin, 10)
	input := make([]byte, 20)

	n, err := ltdreader.Read(input)
	fmt.Printf("LimitReader reads %d bytes: %s\n", n, string(input))
	fmt.Println(n, err)

	discard := make([]byte, 20)
	reader := bufio.NewReader(os.Stdin)
	n, err = reader.Read(discard)
	if err != io.EOF {
		fmt.Printf("Left %d bytes: %s\n", n, string(discard))
	}
}
