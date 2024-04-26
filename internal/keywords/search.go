package keywords

import (
	"context"
	"fmt"
	"jaytaylor.com/html2text"
	"log/slog"
	"net/url"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/cretz/bine/tor"
	"github.com/gocolly/colly/v2"
	"github.com/leapkit/core/envor"
	"github.com/microcosm-cc/bluemonday"
)

// search looks for the specified keyword in the Dark Web.
func (keyword Keyword) Search(service *service) {
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
		source := e.Request.URL.String()

		content = sanitizeHTML(content)
		content = convertContentToText(content)
		content = highlightKeyword(content, keyword.Value)
		if keyword.isContained(content) {
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
			slog.Info(fmt.Sprintf("keyword '%s' was not found in source %s. \n", keyword.Value, source))

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

// isContained returns true when the HTML page contains the keyword
func (keyword Keyword) isContained(content string) bool {
	k := strings.ToLower(keyword.Value)

	genericEmptyResultMessages := []string{
		"nothing found",
		"no results",
	}

	for _, message := range genericEmptyResultMessages {
		if strings.Contains(content, message) {
			return false
		}
	}

	ignoreHTMLElements := ignoredTextElements(k)
	ignoreCoincidences := 0

	for _, element := range ignoreHTMLElements {
		re := regexp.MustCompile(element)
		spot := "[ignore-coincidence]"
		ignoreContent := re.ReplaceAllString(content, spot)
		ignoreCoincidences += strings.Count(ignoreContent, spot)
	}

	numberCoincidences := strings.Count(content, k)
	containsRealResult := numberCoincidences > ignoreCoincidences
	return containsRealResult
}

// sanitizeHTML cleans the html content so it's removed unneeded tags and the escape code.
func sanitizeHTML(contentHTML string) string {
	htmlContent := strings.ToLower(contentHTML)

	p := bluemonday.UGCPolicy()
	p.SkipElementsContent("script")
	cleanedHTML := p.Sanitize(htmlContent)

	return cleanedHTML
}

// convertContentToText converts the HTML to text.
func convertContentToText(content string) string {
	text, err := html2text.FromString(content, html2text.Options{OmitLinks: true, TextOnly: true})
	if err != nil {
		return content
	}

	content = strings.ReplaceAll(text, "\n", " ")

	return content
}

// highlightKeyword applies some styles to highlight the found keyword.
func highlightKeyword(content, keyword string) string {
	content = strings.ToLower(content)
	keyword = strings.ToLower(keyword)

	content = strings.ReplaceAll(content, keyword, fmt.Sprintf(`<span class=f0_5s560nd>%s</span>`, keyword))
	styles := "<div><style>.f0_5s560nd {background-color:#ff8c00b3;font-weight:bold;font-size:16px;}</style>%s</div>"

	return fmt.Sprintf(styles, content)
}

func ignoredTextElements(k string) []string {
	return []string{
		fmt.Sprintf(`<input[^>]*\svalue="%s"[^>]*>`, k),
		fmt.Sprintf(`<nav[^>]*>.*%s.*<\/nav>`, k),
		`search results for\s*“([^”]*)”`,
		fmt.Sprintf(`results for:.*%s`, k),
		fmt.Sprintf(`results for\s.*%s`, k),
		fmt.Sprintf(`result for\s.*%s`, k),
		fmt.Sprintf(`result for.*%s`, k),
		fmt.Sprintf(`%s.*did not match`, k),
	}
}
