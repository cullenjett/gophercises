package link

import (
	"io"
	"strings"

	"golang.org/x/net/html"
)

// Link represents an anchor tag in an HTML document
type Link struct {
	Href string
	Text string
}

// Parse will take in an HTML document and return a
// slice of links parsed from it.
func Parse(r io.Reader) ([]Link, error) {
	doc, err := html.Parse(r)
	if err != nil {
		return nil, err
	}

	nodes := linkNodes(doc)
	var links []Link
	for _, node := range nodes {
		links = append(links, buildLink(node))
	}

	return links, nil
}

func linkNodes(n *html.Node) []*html.Node {
	if n.Type == html.ElementNode && n.Data == "a" {
		return []*html.Node{n}
	}

	var result []*html.Node
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		result = append(result, linkNodes(c)...)
	}
	return result
}

func buildLink(n *html.Node) Link {
	var result Link
	for _, attr := range n.Attr {
		if attr.Key == "href" {
			result.Href = attr.Val
			break
		}
	}
	result.Text = text(n)
	return result
}

func text(n *html.Node) string {
	if n.Type == html.TextNode {
		return n.Data
	}

	if n.Type != html.ElementNode {
		return ""
	}

	var result string
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		result += text(c)
	}
	return strings.Join(strings.Fields(result), " ")
}
