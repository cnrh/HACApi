package utils

import (
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
)

// post posts to a given endpoint with the given formdata, handling failures and returning HTML.
func post(collector *colly.Collector, url, endpoint string, formData map[string]string) (*colly.Collector, *goquery.Selection, error) {
	// Form URL.
	formedUrl := url + endpoint

	collector = collector.Clone()

	// Make a channel to signal if the page is avaliable.
	pageAvaliableChan := make(chan bool, 1)

	// Make a channel to signal any errors.
	errChan := make(chan error, 1)

	// Make a channel to pass HTML through.
	pageHTMLChan := make(chan *colly.HTMLElement, 1)

	// Check if page is avaliable on response.
	collector.OnResponse(func(res *colly.Response) {
		// If final URL is not equal to input URL, the request failed.
		if res.Request.URL.String() != formedUrl {
			pageAvaliableChan <- false
		}
	})

	// Pass HTML to the receiving channel.
	collector.OnHTML("body", func(html *colly.HTMLElement) {
		pageHTMLChan <- html
	})

	// Handle any errors.
	collector.OnError(func(r *colly.Response, err error) {
		errChan <- err
	})

	// Set request headers.
	collector.OnRequest(func(req *colly.Request) {
		req.Headers.Set("Host", strings.Split(url, "//")[1])
		req.Headers.Set("Origin", url)
		req.Headers.Set("Referer", formedUrl)
	})

	// Visit page and wait.
	err := collector.Post(formedUrl, formData)
	collector.Wait()

	// If there was an initial error, return it.
	if err != nil {
		return nil, nil, err
	}

	// Handle errors
	select {
	// Page not avaliable.
	case <-pageAvaliableChan:
		return nil, nil, ErrorPageNotAvaliable
	// Colly error.
	case err := <-errChan:
		return nil, nil, err
	default:
	}

	// Return HTML if it was a success.
	return collector, (<-pageHTMLChan).DOM, nil
}
