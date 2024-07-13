package main

import (
	"fmt"
	"log"
	"time"

	"github.com/gocolly/colly/v2"
)
func main() {
	c := colly.NewCollector(
		colly.Async(true), // Enable asynchronous requests
	)
	// Set rate limits for the scraper
	err := c.Limit(&colly.LimitRule{
		DomainRegexp: `.*\.?facebook\.com`, // Matches facebook.com and any subdomains
		DomainGlob:   "*facebook.com",
		Delay:        2 * time.Second,
		RandomDelay:  1 * time.Second,
		Parallelism:  2,
	})
	if err != nil {
		log.Fatalf("Error setting rate limits: %v", err)
	}
	// On every <a> element with href attribute call callback
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		// Print link
		fmt.Printf("Link found: %q -> %s\n", e.Text, link)
		// Visit link found on page
		// Only those links are visited which are in AllowedDomains
		c.Visit(e.Request.AbsoluteURL(link))
	})
	c.OnError(func(r *colly.Response, err error) {
		fmt.Printf("Request URL: %s failed with response: %d\nError: %v\n", r.Request.URL, r.StatusCode, err)
	})
	// Before making a request print "Visiting ..."
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
	})
	// Start scraping
	err = c.Visit("https://www.facebook.com/")
	if err != nil {
		log.Printf("Unable to reach site: %v", err)
	}
	// Wait until all async tasks are finished
	c.Wait()
}
