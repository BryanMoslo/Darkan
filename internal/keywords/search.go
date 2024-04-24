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

		wg.Done()
	})

	slog.Info(fmt.Sprintf("searching for keyword: '%s'", keyword.Value))
	keywordValue := url.QueryEscape(keyword.Value)

	// List of URLs to scrape
	sources, err := service.SourceList()

	if err != nil {
		slog.Info("error finding source list: ", err)
	}

	for _, source := range sources {
		wg.Add(1)

		u := source.URL + keywordValue

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
