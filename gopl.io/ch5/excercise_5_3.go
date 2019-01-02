// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 123.

// Outline prints the outline of an HTML document tree.
package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"golang.org/x/net/html"
)

//!+
func main() {
	for _, url := range os.Args[1:] {
		if !strings.HasPrefix(url, "http") {
			url = "http://" + url
		}
		DumpTextNodeContents(url)
	}
}

func DumpTextNodeContents(url string) {
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
	dumpTextNodeContents(doc)
	return
}

func dumpTextNodeContents(n *html.Node) {
	switch n.Type {
	case html.ElementNode:
		// Skip unvisiable elements
		if n.Data == "script" || n.Data == "style" {
			return
		}
	case html.TextNode:
		fmt.Printf("%v\n", n.Data)
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		dumpTextNodeContents(c)
	}
}

//!-
