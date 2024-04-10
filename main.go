package main

import (
	"flag"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/gocolly/colly/v2"
)

func main() {
	keyword := flag.String("keyword", "", "Keyword to search for")
	flag.Parse()

	if *keyword == "" {
		fmt.Println("Usage: go run main.go --keyword=\"Acme Inc\"")
		return
	}

	proxyURL, err := url.Parse("socks5://127.0.0.1:9050")
	if err != nil {
		fmt.Println("Error parsing proxy URL:", err)
		return
	}

	c := colly.NewCollector()
	c.SetProxy(proxyURL.String())
	c.SetRequestTimeout(260 * time.Second)

	c.OnHTML("body", func(e *colly.HTMLElement) {
		if strings.Contains(strings.ToLower(e.Text), strings.ToLower(*keyword)) {
			// Getting the whole page for now.
			content, _ := e.DOM.Html()
			fmt.Printf("Keyword '%s' was found in the following HTML content: \n%s\n", *keyword, content)
		} else {
			fmt.Printf("Keyword '%s' was not found. \n", *keyword)
		}
	})

	// libreddit
	onionPage := fmt.Sprintf("http://ecue64yqdxdk3ucrmm2g3irhlvey3wkzcokwi6oodxxwezqk3ak3fhyd.onion/r/popular/search?restrict_sr=on&q=%s", url.QueryEscape(*keyword))

	fmt.Printf("Sniffing '%s' around in the dark web... \n", *keyword)
	err = c.Visit(onionPage)
	if err != nil {
		fmt.Println("Error visiting Onion page:", err)
		return
	}
}
