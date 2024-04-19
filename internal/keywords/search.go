package keywords

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/cretz/bine/tor"
	"github.com/gocolly/colly/v2"
	"github.com/leapkit/core/envor"
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
		if keyword.isContained(e) {
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

			keyword.performCallback(match)
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
		fmt.Sprintf("http://rambleeeqrhty6s5jgefdfdtc6tfgg4jj6svr4jpgk4wjtg3qshwbaad.onion/search?q=%s", url.QueryEscape(keyword.Value)),
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

// performCallback performs the callback to notify that we have the keyword
func (keyword Instance) performCallback(match Match) {
	requestBody, err := json.Marshal(url.Values{
		"Keyword": []string{keyword.Value},
		"Source":  []string{match.Source},
		"Content": []string{match.Content},
	})

	if err != nil {
		slog.Error(fmt.Sprintf("error: %s", err))
		return
	}

	response, err := http.Post(keyword.CallbackURL, "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		slog.Error(fmt.Sprintf("error making request: %s", err))
		return
	}
	defer response.Body.Close()

	slog.Info("performing callback - Response status: " + response.Status)
}

// isContained returns true when the HTML page contains the keyword
func (keyword Instance) isContained(e *colly.HTMLElement) bool {
	htmlContent := strings.ToLower(e.Text)
	k := strings.ToLower(keyword.Value)

	list := []string{
		fmt.Sprintf("couldn’t find any results for “%s”", k),
		fmt.Sprintf("there are no results for %s", k),
		fmt.Sprintf("no results for %s", k),
	}

	for _, item := range list {
		if strings.Contains(htmlContent, item) {
			return false
		}
	}

	return strings.Contains(htmlContent, k)
}
