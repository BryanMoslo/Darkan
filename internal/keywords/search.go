package keywords

import (
	"context"
	"fmt"
	"log/slog"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/cretz/bine/tor"
	"github.com/gocolly/colly/v2"
	"github.com/leapkit/core/envor"
	"golang.org/x/net/html"
)

// search looks for the specified keyword in the Dark Web.
func (keyword Instance) Search(service *service) {
	slog.Info("starting Tor instance")

	t, err := tor.Start(context.TODO(), &tor.StartConf{TempDataDirBase: "tor"})
	if err != nil {
		slog.Error(fmt.Sprintf("failed to start Tor: %s", err.Error()))
		return
	}

	defer t.Close()

	c := colly.NewCollector(
		colly.MaxDepth(1),
		colly.Async(true),
	)

	torProxy := envor.Get("TOR_PROXY", "socks5://127.0.0.1:9050")
	if err := c.SetProxy(torProxy); err != nil {
		slog.Error(fmt.Sprintf("error setting up a proxy: %s", err.Error()))
		return
	}

	c.SetRequestTimeout(5 * time.Minute)

	c.Limit(&colly.LimitRule{
		Parallelism: 2,
	})

	var wg sync.WaitGroup

	c.OnHTML("body", func(e *colly.HTMLElement) {
		content, _ := e.DOM.Html()
		if keyword.isContained(content) {
			source := e.Request.URL.String()
			slog.Info(fmt.Sprintf("keyword '%s' was found in source: '%s'", keyword.Value, source))

			match := Match{
				KeywordID: keyword.ID,
				Content:   content,
				Source:    source,
				FoundAt:   time.Now(),
			}

			err := service.CreateMatch(&match)
			if isDuplicateKeyError(err) {
				slog.Error("match already exists for keyword and source URL")
				return
			}

			if err != nil {
				slog.Error(fmt.Sprintf("failed to create a match for keyword '%s' in source: %s, error: %s", keyword.Value, source, err.Error()))
			}

			// TODO:
			// Make a POST request to the keyword.CallbackURL saying we have found the keyword.
		} else {
			slog.Info(fmt.Sprintf("keyword '%s' was not found. \n", keyword.Value))

			// TODO:
			// Storing some info about the research we've done (?)
			// Keep the keyword as Found = false so once this run again (via CLI/CronJob), it will look for all keywords that have not been found yet (until?: undefined for now)
		}

		wg.Done()
	})

	c.OnError(func(_ *colly.Response, err error) {
		slog.Info(fmt.Sprintf("something went wrong: %s", err))
	})

	// TODO:
	// Add More sources

	// List of URLs to scrape
	urls := []string{
		// fmt.Sprintf("http://ecue64yqdxdk3ucrmm2g3irhlvey3wkzcokwi6oodxxwezqk3ak3fhyd.onion/r/popular/search?restrict_sr=on&q=%s", url.QueryEscape(keyword.Value)), UNAVAILABLE
		fmt.Sprintf("https://www.reddittorjg6rue252oqsxryoxengawnmo46qy4kyii5wtqnwfj4ooad.onion/search?q=%s", url.QueryEscape(keyword.Value)),
		fmt.Sprintf("http://rambleeeqrhty6s5jgefdfdtc6tfgg4jj6svr4jpgk4wjtg3qshwbaad.onion/search?q=%s", url.QueryEscape(keyword.Value)), // Returns false positive
		fmt.Sprintf("https://www.bbcnewsd73hkzno2ini43t4gblxvycyac5aw4gnv7t2rccijh7745uqd.onion/search?q=%s", url.QueryEscape(keyword.Value)),
	}

	slog.Info(fmt.Sprintf("searching for keyword: '%s'", keyword.Value))
	for _, u := range urls {
		wg.Add(1)

		err := c.Visit(u)
		if err != nil {
			slog.Error(fmt.Sprintf("error visiting %s: %s", u, err))
			return
		}
	}

	wg.Wait()
	c.Wait()

	slog.Info(fmt.Sprintf("search completed for: '%s'", keyword.Value))
}

// isContained returns true when the given text contains the keyword but not in an URL
func (keyword Instance) isContained(text string) bool {
	k := strings.ToLower(keyword.Value)

	doc, err := html.Parse(strings.NewReader(text))
	if err != nil {
		slog.Error(fmt.Sprintf("error by trying to parse: %s", err))
		return false
	}

	var found bool
	var f func(*html.Node)
	f = func(n *html.Node) {
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			data := strings.ToLower(c.Data)

			_, err := url.ParseRequestURI(data)
			isURL := err == nil

			if !isURL && strings.Contains(data, k) {
				found = true
				return
			}

			f(c)
		}
	}

	f(doc)
	return found
}
