package main

import (
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"golang.org/x/net/html"
)

func main() {
	urlFlag := flag.String("url", "https://gophercises.com", "the root URL used to build a sitemap")
	maxDepth := flag.Int("depth", 3, "maximum number of links deep to traverse")
	flag.Parse()

	pages := bfs(*urlFlag, *maxDepth)
	for _, page := range pages {
		fmt.Println(page)
	}

	// links := getLinks(*urlFlag)
	// for _, link := range links {
	// 	fmt.Println(link)
	// }
}

func bfs(urlString string, maxDepth int) []string {
	// apparently a struct{} uses less memory than a bool
	visited := make(map[string]struct{})
	var q map[string]struct{}
	nextQ := map[string]struct{}{
		urlString: struct{}{},
	}

	for i := 0; i <= maxDepth; i++ {
		q, nextQ = nextQ, make(map[string]struct{})
		for url := range q {
			if _, ok := visited[url]; ok {
				continue
			}
			visited[url] = struct{}{}
			for _, link := range getLinks(url) {
				nextQ[link] = struct{}{}
			}
		}
	}

	result := make([]string, 0, len(visited)) // minor optimization vs creating zero value []string
	for u := range visited {
		result = append(result, u)
	}
	return result
}

func getLinks(urlString string) []string {
	res, err := http.Get(urlString)
	if err != nil {
		return []string{}
	}
	defer res.Body.Close()

	reqURL := res.Request.URL
	baseURL := &url.URL{
		Scheme: reqURL.Scheme,
		Host:   reqURL.Host,
	}
	base := baseURL.String()
	hrefs := getHrefs(res.Body, base)

	return filter(hrefs, func(link string) bool {
		return strings.HasPrefix(link, base)
	})
}

func getHrefs(r io.Reader, base string) []string {
	links, err := ParseLinks(r)
	if err != nil {
		return []string{}
	}

	var hrefs []string
	for _, l := range links {
		switch {
		case strings.HasPrefix(l.Href, "/"):
			hrefs = append(hrefs, base+l.Href)
		case strings.HasPrefix(l.Href, "http"):
			hrefs = append(hrefs, l.Href)
		}
	}

	return hrefs
}

func filter(items []string, keepFn func(string) bool) []string {
	var result []string

	for _, item := range items {
		if keepFn(item) {
			result = append(result, item)
		}
	}

	return result
}

// Link represents an anchor tag in an HTML document
type Link struct {
	Href string
	Text string
}

// ParseLinks will take in an HTML document and return a
// slice of links parsed from it.
func ParseLinks(r io.Reader) ([]Link, error) {
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
