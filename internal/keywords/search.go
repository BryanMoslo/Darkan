package keywords

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/cretz/bine/tor"
	"github.com/gocolly/colly/v2"
	"github.com/leapkit/core/envor"
	"github.com/microcosm-cc/bluemonday"
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
		if keyword.isContained(sanitizeHTML(content)) {
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
func (keyword Instance) isContained(content string) bool {
	k := strings.ToLower(keyword.Value)

	genericEmptyResultMessages := []string{
		"nothing found",
		"no results",
	}

	for _, message := range genericEmptyResultMessages {
		if strings.Contains(content, message) {
			fmt.Println("------->", false)
			return false
		}
	}

	ignoreHTMLElements := []string{
		fmt.Sprintf(`<input[^>]*\svalue="%s"[^>]*>`, k),
		fmt.Sprintf(`<nav[^>]*>.*%s.*<\/nav>`, k),
		`search results for\s*“([^”]*)”`,
		fmt.Sprintf(`results for:.*%s`, k),
		fmt.Sprintf(`results for\s.*%s`, k),
		fmt.Sprintf(`result for\s.*%s`, k),
	}

	ignoreCoincidences := 0
	for _, element := range ignoreHTMLElements {
		re := regexp.MustCompile(element)
		spot := "[ignore-coincidence]"
		ignoreContent := re.ReplaceAllString(content, spot)
		ignoreCoincidences += strings.Count(ignoreContent, spot)
	}

	numberCoincidences := strings.Count(content, k)
	containsRealResult := numberCoincidences > ignoreCoincidences
	// fmt.Println("------->", numberCoincidences)
	// fmt.Println("------->", ignoreCoincidences)
	// fmt.Println("------->", containsRealResult)
	return containsRealResult
}

// sanitizeHTML cleans the html content so it's removed
func sanitizeHTML(contentHTML string) string {
	htmlContent := strings.ToLower(contentHTML)

	p := bluemonday.UGCPolicy()
	cleanedHTML := p.Sanitize(htmlContent)

	cleanedHTML = strings.ReplaceAll(cleanedHTML, "\n", "")
	doc, err := html.Parse(strings.NewReader(cleanedHTML))

	if err != nil {
		log.Fatal(err)
		return ""
	}

	removeScript(doc)

	buf := bytes.NewBuffer([]byte{})
	if err := html.Render(buf, doc); err != nil {
		log.Fatal(err)
		return ""
	}

	return buf.String()

}

// removeScript removes the script tags content in an HTML node.
func removeScript(n *html.Node) {
	if n.Type == html.ElementNode && n.Data == "script" {
		n.Parent.RemoveChild(n)
		return // script tag is gone...
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		removeScript(c)
	}
}
