// main.go
package main

import (
	"darkan/search"
	"flag"
	"fmt"
)

func main() {
	keyword := flag.String("keyword", "", "Keyword to search for")
	flag.Parse()

	if *keyword == "" {
		fmt.Println("Usage: go run main.go --keyword=\"Acme Inc\"")
		return
	}

	search.SearchDarkWeb(*keyword)
}
