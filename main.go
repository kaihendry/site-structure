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
func crawl(url string, ch chan string, chFinished chan bool) {

	var b io.ReadCloser
	defer func() {
		// Notify that we're done after this function
		chFinished <- true
	}()

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

	log.Println("tokenizing")

	for {
		tt := z.Next()

		log.Println("walking", tt)
		switch {
		case tt == html.ErrorToken:
			log.Println("End of the document, we're done")
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

			ch <- url
		}

	}
}

func main() {
	foundUrls := make(map[string]bool)
	seedUrls := os.Args[1:]

	// Channels
	chUrls := make(chan string)
	chFinished := make(chan bool)

	// Kick off the crawl process (concurrently)
	for _, url := range seedUrls {
		go crawl(url, chUrls, chFinished)
	}

	// Subscribe to both channels
	for c := 0; c < len(seedUrls); {
		select {
		case url := <-chUrls:
			log.Println("Found", url)
			foundUrls[url] = true
		case <-chFinished:
			c++
		}
	}

	// We're done! Print the results...

	fmt.Println("\nFound", len(foundUrls), "unique urls:")

	for url, _ := range foundUrls {
		fmt.Println(" - " + url)
	}

	close(chUrls)
}
