package crawler

import (
	"log"
	"sync"

	"github.com/spike01/webcrawler/client"
	"github.com/spike01/webcrawler/page"
)

type Crawler struct {
	client  *client.Client
	verbose bool
	pages   []*page.Page

	mutex   sync.Mutex
	visited map[string]bool
}

func NewCrawler(c *client.Client, verbose bool) *Crawler {
	return &Crawler{
		client:  c,
		verbose: verbose,
		pages:   make([]*page.Page, 0),
		visited: make(map[string]bool),
	}
}

func (c *Crawler) Crawl(url string) ([]*page.Page, error) {
	c.crawlRecursive(url)
	return c.pages, nil
}

func (c *Crawler) crawlRecursive(url string) {
	var wg sync.WaitGroup

	v := c.visit(url)
	if v {
		return
	}

	page := c.getPage(url)
	c.pages = append(c.pages, page)

	for _, link := range page.Links {
		wg.Add(1)
		go func(link string) {
			defer wg.Done()
			c.crawlRecursive(link)
		}(link.String())
	}
	wg.Wait()
	return
}

func (c *Crawler) visit(url string) bool {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	_, ok := c.visited[url]
	if ok {
		return true
	}
	c.visited[url] = true

	return false
}

func (c *Crawler) getPage(url string) *page.Page {
	if c.verbose {
		log.Printf("Fetching: %s", url)
	}

	body, err := c.client.Get(url)
	if err != nil {
		log.Printf("Unable to fetch url: %s. Reason: %s", url, err)
	}

	page, err := page.NewPage(url, body)
	if err != nil {
		log.Printf("Unable to create page from url: %s. Reason: %s", url, err)
	}
	return page
}
