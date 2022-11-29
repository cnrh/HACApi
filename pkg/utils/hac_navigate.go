package utils

import (
	"errors"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
)

// NavigateTo navigates a collector to a specified URL, handling failures and returning HTML.
func NavigateTo(baseCollector *colly.Collector, base, url string) (*colly.Collector, *goquery.Selection, error) {
	// Form full URL
	formedUrl := "https://" + base + url

	// Make a copy of the collector
	collector := baseCollector.Clone()

	// Channel to confirm if the page was avaliable
	pageAvaliableChan := make(chan bool, 1)

	// Channel to pass HTML through
	pageHTMLChan := make(chan *colly.HTMLElement, 1)

	// Check if page is avaliable on response
	collector.OnResponse(func(res *colly.Response) {
		// If final URL equal to input URL, it's successful
		pageAvaliableChan <- (res.Request.URL.String() == formedUrl)
	})

	// Pass HTML to the receiving channel
	collector.OnHTML("body", func(html *colly.HTMLElement) {
		pageHTMLChan <- html
	})

	// Visit page
	err := collector.Visit(formedUrl)
	collector.Wait()

	// Return false if not avaliable, or if there was an error
	if err != nil {
		return nil, nil, err
	}

	if !(<-pageAvaliableChan) {
		return nil, nil, errors.New("page not avaliable")
	}

	// Return HTML if all is well
	return collector, (<-pageHTMLChan).DOM, nil
}
