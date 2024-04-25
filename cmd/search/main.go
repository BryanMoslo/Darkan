package main

import (
	"darkan/internal"
	"darkan/internal/keywords"
	"fmt"
	"log/slog"
	"sync"
)

func main() {
	conn, err := internal.DB()
	if err != nil {
		slog.Error(err.Error())
		return
	}

	keywordService := keywords.NewService(conn)
	keywordsToSearch, err := keywordService.UnfoundKeywords()
	if err != nil {
		slog.Error(err.Error())
		return
	}

	slog.Info(fmt.Sprintf("%v keyword(s) found to search in the dark web", len(keywordsToSearch)))

	var wg sync.WaitGroup
	for _, keyword := range keywordsToSearch {
		wg.Add(1)
		go func(k keywords.Keyword) {
			defer wg.Done()
			k.Search(keywordService)
		}(keyword)
	}

	wg.Wait()
	slog.Info("[done] process completed for all keywords")
}
