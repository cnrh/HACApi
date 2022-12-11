package utils

import (
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
)

// post posts to a given endpoint with the given formdata, handling failures and returning HTML.
func post(collector *colly.Collector, url, endpoint string, formData map[string]string) (*colly.Collector, *goquery.Selection, error) {
	// Form URL
	formedUrl := url + endpoint

	// Make a copy of the collector
	collector = collector.Clone()

	// Channel to confirm if the page was avaliable
	pageAvaliableChan := make(chan bool, 1)

	// Channel to handle any errors
	errChan := make(chan error, 1)

	// Channel to pass HTML through
	pageHTMLChan := make(chan *colly.HTMLElement, 1)

	// Check if page is avaliable on response
	collector.OnResponse(func(res *colly.Response) {
		// If final URL equal to input URL, it's successful
		if res.Request.URL.String() != formedUrl {
			pageAvaliableChan <- false
		}
	})

	// Pass HTML to the receiving channel
	collector.OnHTML("body", func(html *colly.HTMLElement) {
		pageHTMLChan <- html
	})

	// Handle any errors
	collector.OnError(func(r *colly.Response, err error) {
		errChan <- err
	})

	collector.OnRequest(func(req *colly.Request) {
		req.Headers.Set("Host", strings.Split(url, "//")[1])
		req.Headers.Set("Origin", url)
		req.Headers.Set("Referer", formedUrl)
	})

	// Visit page
	err := collector.Post(formedUrl, formData)
	collector.Wait()

	// Return false if not avaliable, or if there was an error
	if err != nil {
		return nil, nil, err
	}

	// Handle any other errors
	select {
	case <-pageAvaliableChan:
		return nil, nil, ErrorPageNotAvaliable
	case err := <-errChan:
		return nil, nil, err
	default:
	}

	// Return HTML if all is well
	return collector, (<-pageHTMLChan).DOM, nil
}
