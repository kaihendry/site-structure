package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"strings"

	"golang.org/x/net/html"
)

var basePath string

// Helper function to pull the href attribute from a Token
func getHref(t html.Token) (ok bool, href string) {
	// Iterate over token attributes until we find an "href"
	for _, a := range t.Attr {
		if a.Key == "href" {
			href = a.Val
			ok = true
		}
	}

	// "bare" return will return the variables (ok, href) as
	// defined in the function definition
	return
}

// Extract links from a Webpage
func crawl(url string) (urls []string, err error) {

	var b io.ReadCloser

	if strings.HasPrefix(url, "http") {
		log.Printf("Fetching: %s", url)
		resp, err := http.Get(url)
		if err != nil {
			log.Println("Failed to fetch", url)
			return urls, err
		}
		b = resp.Body
		defer b.Close()
	} else {
		log.Printf("Opening: %s", url)
		b, err = os.Open(url)
		if err != nil {
			return urls, err
		}
	}

	z := html.NewTokenizer(b)

	for {
		tt := z.Next()
		switch {
		case tt == html.ErrorToken:
			return
		case tt == html.StartTagToken:
			t := z.Token()

			// Check if the token is an <a> tag
			isAnchor := t.Data == "a"
			if !isAnchor {
				continue
			}

			// Extract the href value, if there is one
			ok, url := getHref(t)
			if !ok {
				continue
			}
			urls = append(urls, url)
		}
	}
}

func main() {
	flag.Parse()
	root := flag.Args()[0]
	basePath = path.Dir(root)
	root = path.Base(root)
	log.Println("Starting with", basePath, root)
	spider([]string{root}, "")
}

func spider(urls []string, parent string) {
	log.Printf("Spidering: %v\n", urls)
	for _, u := range urls {
		log.Println("Crawling", path.Join(basePath, u))
		foundUrls, err := crawl(path.Join(basePath, u))
		if err != nil {
			log.Printf("Ignoring %v", err)
			return
		}
		if parent != "" {
			fmt.Printf("%s,%s\n", u, parent)
		}
		spider(foundUrls, u)
	}
}
