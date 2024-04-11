package search

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/cretz/bine/tor"
	"github.com/gocolly/colly/v2"
)

// SearchDarkWeb searches for the specified keyword in the Dark Web.
func SearchDarkWeb(keyword string) {
	fmt.Println("Starting and registering Tor hidden service...")
	t, err := tor.Start(context.TODO(), &tor.StartConf{})
	if err != nil {
		fmt.Printf("Failed to start Tor: %v\n", err)
		return
	}
	defer t.Close()

	c := colly.NewCollector(
		colly.MaxDepth(1),
		colly.Async(true),
	)

	torProxy := os.Getenv("TOR_PROXY")
	if err := c.SetProxy(torProxy); err != nil {
		fmt.Println("Error setting up a proxy:", err)
		return
	}

	c.SetRequestTimeout(5 * time.Minute)

	c.OnHTML("body", func(e *colly.HTMLElement) {
		if strings.Contains(strings.ToLower(e.Text), strings.ToLower(keyword)) {
			// Getting the whole page for now.
			content, _ := e.DOM.Html()
			fmt.Printf("Keyword '%s' was found in the following HTML content: \n%s\n", keyword, content)
		} else {
			fmt.Printf("Keyword '%s' was not found. \n", keyword)
		}
	})

	// libreddit
	onionPage := fmt.Sprintf("http://ecue64yqdxdk3ucrmm2g3irhlvey3wkzcokwi6oodxxwezqk3ak3fhyd.onion/r/popular/search?restrict_sr=on&q=%s", url.QueryEscape(keyword))

	fmt.Printf("Sniffing '%s' around in the dark web... \n", keyword)
	err = c.Visit(onionPage)
	if err != nil {
		fmt.Println("Error visiting Onion page:", err)
		return
	}

	c.Wait()
}
