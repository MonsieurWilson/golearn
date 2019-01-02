package main

import (
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
		images, err := ElementsByTagName(url, "img")
		// headings, err := ElementsByTagName(url, "h1", "h2", "h3", "h4")
		if err != nil {
			fmt.Errorf("ElementsByTagName error: %v\n", err)
		} else {
			fmt.Println("Images:")
			for _, n := range images {
				fmt.Printf("<%s", n.Data)
				for _, attr := range n.Attr {
					fmt.Printf(" %s=%s", attr.Key, attr.Val)
				}
				fmt.Printf(">\n")
			}
		}
	}
}

func ElementsByTagName(url string, names ...string) (n []*html.Node, err error) {
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
	n = elementsByTagName(doc, names...)
	return
}

func elementsByTagName(n *html.Node, names ...string) (res []*html.Node) {
	if n.Type == html.ElementNode {
		for _, val := range names {
			if n.Data == val {
				res = append(res, n)
			}
		}
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		res = append(res, elementsByTagName(c, names...)...)
	}

	return
}
