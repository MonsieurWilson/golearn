package main

import (
	"bufio"
	"bytes"
	"fmt"
	"net/http"
	"os"
	"strings"

	"golang.org/x/net/html"
)

func main() {
	for _, url := range os.Args[1:] {
		if !strings.HasPrefix(url, "http") {
			url = "http://" + url
		}
		nwords, nimages, err := CountWordsAndImages(url)
		if err != nil {
			fmt.Errorf("CountWordsAndImages error: %v\n", err)
		} else {
			fmt.Printf("Words: %v, images: %v\n", nwords, nimages)
		}
	}
}

// CountWordsAndImages does an HTTP GET request for the HTML
// document url and returns the number of words and images in it.
func CountWordsAndImages(url string) (words, images int, err error) {
	resp, err := http.Get(url)
	if err != nil {
		return
	}
	doc, err := html.Parse(resp.Body)
	resp.Body.Close()
	if err != nil {
		err = fmt.Errorf("parsing HTML error: %s", err)
		return
	}
	words, images = countWordsAndImages(doc)
	return
}

func countWordsAndImages(n *html.Node) (words, images int) {
	// count the number of words and images
	switch n.Type {
	case html.ElementNode:
		if n.Data == "script" || n.Data == "style" {
			return
		}
		if n.Data == "img" {
			images++
		}
	case html.TextNode:
		buf := bytes.NewBufferString(n.Data)
		input := bufio.NewScanner(buf)
		input.Split(bufio.ScanWords)

		for input.Scan() {
			words++
		}
	}

	// count recursively for child nodes
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		nwords, nimages := countWordsAndImages(c)
		words, images = words+nwords, images+nimages
	}

	return
}
