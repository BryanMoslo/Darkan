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
	keywordsToSearch, err := keywordService.All()
	if err != nil {
		slog.Error(err.Error())
		return
	}

	slog.Info(fmt.Sprintf("%v keyword(s) to search in the dark web.", len(keywordsToSearch)))

	var wg sync.WaitGroup
	for _, keyword := range keywordsToSearch {
		wg.Add(1)
		go func(k keywords.Instance) {
			defer wg.Done()
			k.Search()
		}(keyword)
	}

	wg.Wait()
	slog.Info("All keywords searched.")
}
