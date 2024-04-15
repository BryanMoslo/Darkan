package keywords

import (
	"context"
	"fmt"
	"log/slog"
	"net/url"
	"strings"
	"time"

	"github.com/cretz/bine/tor"
	"github.com/gocolly/colly/v2"
	"github.com/leapkit/core/envor"
)

// SearchDarkWeb searches for the specified keyword in the Dark Web.
func (keyword Instance) search() error {
	slog.Info("Starting Tor instance...")

	t, err := tor.Start(context.TODO(), &tor.StartConf{})
	if err != nil {
		slog.Info(fmt.Sprintf("failed to start Tor: %s", err.Error()))
		return fmt.Errorf("failed to start Tor: %v", err)
	}

	defer t.Close()

	c := colly.NewCollector(
		colly.MaxDepth(1),
		colly.Async(true),
	)

	torProxy := envor.Get("TOR_PROXY", "socks5://127.0.0.1:9050")
	if err := c.SetProxy(torProxy); err != nil {
		slog.Info(fmt.Sprintf("error setting up a proxy: %s", err.Error()))
		return fmt.Errorf("error setting up a proxy: %v", err)
	}

	c.SetRequestTimeout(5 * time.Minute)

	c.OnHTML("body", func(e *colly.HTMLElement) {
		if strings.Contains(strings.ToLower(e.Text), strings.ToLower(keyword.Value)) {
			content, _ := e.DOM.Html()
			slog.Info(fmt.Sprintf("Keyword '%s' was found in the following HTML content: \n%s\n", keyword.Value, content))

			// TODO:
			// Save the source's info (URL and Content).
			// Update the Keyword instance to set keyword.Found = true
			// Make a POST request to the keyword.CallbackURL saying we have found the keyword.
		} else {
			slog.Info(fmt.Sprintf("Keyword '%s' was not found. \n", keyword.Value))

			// TODO:
			// Storing some info about the research we've done (?)
			// Keep the keyword as Found = false so once this run again (via CLI/CronJob), it will look for all keywords that have not been found yet (until?: undefined for now)
		}
	})

	// TODO:
	// Add More sources
	// Implement an efficient way to perform this search (concurreny?)

	// libreddit
	onionPage := fmt.Sprintf("http://ecue64yqdxdk3ucrmm2g3irhlvey3wkzcokwi6oodxxwezqk3ak3fhyd.onion/r/popular/search?restrict_sr=on&q=%s", url.QueryEscape(keyword.Value))

	fmt.Printf("Sniffing '%s' around in the dark web... \n", keyword.Value)
	err = c.Visit(onionPage)
	if err != nil {
		return fmt.Errorf("error visiting Onion page: %v", err)
	}

	c.Wait()

	return nil
}
