package keywords

type KeywordMatch struct {
	Instance
	Match
}

type KeywordsMatches []KeywordMatch

func (km KeywordMatch) PerformCallback() {
	km.performCallback(km.Match)
}
