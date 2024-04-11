// main.go
package main

import (
	"darkan/search"
	"flag"
	"fmt"
	"log"

	"github.com/joho/godotenv"
)

func main() {
	keyword := flag.String("keyword", "", "Keyword to search for")
	flag.Parse()

	if *keyword == "" {
		fmt.Println("Usage: go run main.go --keyword=\"Acme Inc\"")
		return
	}

	if err := godotenv.Load(".env"); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	search.SearchDarkWeb(*keyword)
}
