package utils

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
)

// NavigateCollectorTo navigates a collector to a specified URL, checking if it succeeds.
func NavigateTo(baseCollector *colly.Collector, url string) (*colly.Collector, *goquery.Selection, bool) {
	// Make a copy of the collector
	collector := baseCollector.Clone()

	// Channel to confirm if the page was avaliable
	pageAvaliableChan := make(chan bool, 1)

	// Channel to pass HTML through
	pageHTMLChan := make(chan *colly.HTMLElement, 1)

	// Check if page is avaliable on response
	collector.OnResponse(func(res *colly.Response) {
		// If final URL equal to input URL, it's successful
		pageAvaliableChan <- (res.Request.URL.String() == url)
	})

	// Pass HTML to the receiving channel
	collector.OnHTML("body", func(html *colly.HTMLElement) {
		pageHTMLChan <- html
	})

	// Visit page
	err := collector.Visit(url)
	collector.Wait()

	// Return false if not avaliable, or if there was an error
	if isPageAvaliable := <-pageAvaliableChan; err != nil || !isPageAvaliable {
		return nil, nil, false
	}

	// Return HTML if all is well
	return collector, (<-pageHTMLChan).DOM, true
}
