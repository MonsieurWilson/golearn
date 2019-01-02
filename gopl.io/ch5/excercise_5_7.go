// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 133.

// Outline prints the outline of an HTML document tree.
package main

import (
	"bytes"
	"fmt"
	"net/http"
	"os"
	"strings"

	"golang.org/x/net/html"
)

func main() {
	for _, url := range os.Args[1:] {
		outline(url)
	}
}

func outline(url string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	doc, err := html.Parse(resp.Body)
	if err != nil {
		return err
	}

	//!+call
	forEachNode(doc, startElement, endElement)
	//!-call

	return nil
}

//!+forEachNode
// forEachNode calls the functions pre(x) and post(x) for each node
// x in the tree rooted at n. Both functions are optional.
// pre is called before the children are visited (preorder) and
// post is called after (postorder).
func forEachNode(n *html.Node, pre, post func(n *html.Node)) {
	if pre != nil {
		pre(n)
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		forEachNode(c, pre, post)
	}

	if post != nil {
		post(n)
	}
}

//!-forEachNode

//!+startend
var depth int

func startElement(n *html.Node) {
	switch n.Type {
	case html.ElementNode:
		var buf bytes.Buffer
		buf.WriteString(fmt.Sprintf("%*s<%s", depth*2, "", n.Data))
		for _, attr := range n.Attr {
			buf.WriteString(fmt.Sprintf(" %s=%s", attr.Key, attr.Val))
		}
		if c := n.FirstChild; c == nil {
			buf.WriteString("/>")
		} else {
			buf.WriteString(">")
		}
		fmt.Println(buf.String())
		depth++
	case html.TextNode:
		fmt.Printf("%*s%s\n", depth*2, "", strings.Trim(n.Data, " \n"))
	}
}

func endElement(n *html.Node) {
	if n.Type == html.ElementNode {
		depth--
		if c := n.FirstChild; c != nil {
			fmt.Printf("%*s</%s>\n", depth*2, "", n.Data)
		}
	}
}

//!-startend
