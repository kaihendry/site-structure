package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
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
func crawl(url string) (urls []string) {

	var b io.ReadCloser

	if strings.HasPrefix(url, "http") {
		log.Printf("Fetching: %s", url)
		resp, err := http.Get(url)

		if err != nil {
			fmt.Println("ERROR: Failed to crawl:", url)
			return
		}
		b = resp.Body
		defer b.Close()
	} else {
		log.Printf("Opening: %s", url)
		var err error
		b, err = os.Open(url)
		if err != nil {
			log.Fatal(err)
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
	foundUrls := crawl(os.Args[1])
	fmt.Println("\nFound", len(foundUrls), "unique urls:")
	for _, u := range foundUrls {
		fmt.Println("Open", u)
		urls := crawl("test/" + u)
		fmt.Println("\nFound", len(urls), "unique urls:")
		for _, v := range urls {
			fmt.Printf("%s points to %s", u, v)
		}
	}

}
