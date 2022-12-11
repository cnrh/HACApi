package utils

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
)

// Scraper is the struct used to inject the HAC scraping
// dependency into the Server struct.
type Scraper struct{}

func (scraper Scraper) Login(url, username, password string) (*colly.Collector, error) {
	return login(url, username, password)
}

func (scraper Scraper) Navigate(collector *colly.Collector, url, endpoint string) (*colly.Collector, *goquery.Selection, error) {
	return navigate(collector, url, endpoint)
}

func (scraper Scraper) Post(collector *colly.Collector, url, endpoint string, formData map[string]string) (*colly.Collector, *goquery.Selection, error) {
	return post(collector, url, endpoint, formData)
}

func NewScraper() *Scraper {
	return &Scraper{}
}
