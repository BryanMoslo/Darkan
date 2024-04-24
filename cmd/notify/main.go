package main

import (
	"darkan/internal"
	"darkan/internal/keywords"
	"log/slog"
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

	for _, keywordMatch := range keywordMatchList {
		<-taskLimiter.C
		keywordMatch.PerformCallback()
	}

	slog.Info("[done] notification process completed.")
}
