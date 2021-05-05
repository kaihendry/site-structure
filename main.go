package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"strings"

	"golang.org/x/net/html"
)

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
	foundUrls, err := crawl(os.Args[1])
	dirPath := path.Dir(os.Args[1])
	basePath := path.Base(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	for _, u := range foundUrls {
		log.Println("Crawling", u)
		urls, err := crawl(path.Join(dirPath, u))
		if err != nil {
			log.Printf("Ignoring %v", err)
			break
		}
		fmt.Printf("%s,%s\n", u, basePath)
		for _, v := range urls {
			fmt.Println(v)
		}
	}

}
