package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"

	"golang.org/x/net/html"
)

var seen map[string]bool

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

	if seen[url] {
		log.Println("seen", url, "before, not crawling")
		return
	}

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

	log.Println("marking", url, "as seen")
	seen[url] = true

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
	seen = make(map[string]bool)
	root := os.Args[1]
	if root == "" {
		log.Fatal("Missing argument")
	}
	spider([]string{root}, "")
}

func spider(urls []string, parent string) {
	log.Printf("Spidering: %v\n", urls)
	for _, u := range urls {
		if !strings.HasPrefix(u, "http") {
			basePath := path.Dir(u)
			if basePath != "" {
				os.Chdir(basePath)
			}
			u = path.Base(u)
			cwd, err := os.Getwd()
			if err != nil {
				log.Fatal(err)
			}
			log.Println("Local crawl", path.Join(cwd, u))
		}

		foundUrls, err := crawl(u)
		if err != nil {
			log.Printf("Ignoring %v", err)
			return
		}
		if parent != "" {
			fmt.Printf("%s,%s\n", u, parent)
		}
		log.Println("Found URLs", foundUrls)

		if strings.HasPrefix(u, "http") {
			var resolvedURLs []string
			for _, furl := range foundUrls {
				f, err := url.Parse(furl)
				if err != nil {
					log.Fatal(err)
				}
				base, err := url.Parse(u)
				if err != nil {
					log.Fatal(err)
				}
				resolvedURLs = append(resolvedURLs, base.ResolveReference(f).String())
			}
			log.Println("http urls to spider", resolvedURLs)
			spider(resolvedURLs, u)
		} else {
			spider(foundUrls, u)
		}
	}
}
