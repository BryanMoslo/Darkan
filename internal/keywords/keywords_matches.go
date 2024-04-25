package keywords

import (
	"net/url"
)

type KeywordMatch struct {
	Keyword
	Match
}

type KeywordsMatches []KeywordMatch

func (km KeywordsMatches) MatchesByClient() (clientMatches map[string][]url.Values) {
	clientMatches = map[string][]url.Values{}

	for _, keywordMatch := range km {
		clientMatches[keywordMatch.CallbackURL] = append(clientMatches[keywordMatch.CallbackURL], url.Values{
			"Keyword": []string{keywordMatch.Value},
			"Source":  []string{keywordMatch.Source},
			"Content": []string{keywordMatch.Content},
		})
	}

	return clientMatches
}
