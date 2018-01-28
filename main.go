package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/spike01/webcrawler/client"
	"github.com/spike01/webcrawler/crawler"
)

var verbose = flag.Bool("v", false, "make the operation more talkative")

func main() {
	flag.Parse()

	if len(flag.Args()) != 1 {
		fmt.Println("try 'webcrawler --help' for more information")
		os.Exit(1)
	}

	domain := flag.Args()[0]

	client := client.NewClient(&http.Client{
		Timeout: time.Second * 5,
	})

	crawler := crawler.NewCrawler(client, *verbose)

	pages, err := crawler.Crawl(domain)
	if err != nil {
		log.Fatalf("Could not crawl URL: %s", domain)
	}

	if *verbose {
		log.Printf("Retrieved %d pages", len(pages))
	}

	for _, page := range pages {
		fmt.Println(page)
	}

	os.Exit(0)
}
