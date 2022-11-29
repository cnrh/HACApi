package utils

import (
	"errors"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
)

// PostTo posts to a given endpoint with the given formdata, handling failures and returning HTML.
func PostTo(baseCollector *colly.Collector, base, url string, formData map[string]string) (*colly.Collector, *goquery.Selection, error) {
	// Form URL
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

	collector.OnRequest(func(req *colly.Request) {
		req.Headers.Set("Host", base)
		req.Headers.Set("Origin", "https://"+base)
		req.Headers.Set("Referer", formedUrl)
	})

	// Visit page
	err := collector.Post(formedUrl, formData)
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
