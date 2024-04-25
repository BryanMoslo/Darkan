package main

import (
	"bytes"
	"darkan/internal"
	"darkan/internal/keywords"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"time"
)

const taskCreationRate = 400

var taskLimiter = time.NewTicker(time.Duration(20) * time.Second / taskCreationRate)

func main() {
	conn, err := internal.DB()
	if err != nil {
		slog.Error(err.Error())
		return
	}
	keywordService := keywords.NewService(conn)
	keywordMatchList, err := keywordService.KeywordMatchList()
	matchByClient := keywordMatchList.MatchesByClient()

	for callbackURL, body := range matchByClient {
		<-taskLimiter.C
		performCallback(callbackURL, body)
	}

	slog.Info("[done] notification process completed.")
}

// performCallback performs the callback to notify that we have the keyword
func performCallback(url string, matches []url.Values) {
	requestBody, err := json.Marshal(matches)

	if err != nil {
		slog.Error(fmt.Sprintf("error: %s", err))
		return
	}

	response, err := http.Post(url, "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		slog.Error(fmt.Sprintf("error making request: %s", err))
		return
	}
	defer response.Body.Close()

	slog.Info("performing callback - Response status: " + response.Status)
}
