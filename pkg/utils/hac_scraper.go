package utils

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
)

// Scraper is the struct used to inject the HAC scraping
// dependency into the Server struct.
type Scraper struct{}

func (scraper Scraper) Login(base, username, password string) (*colly.Collector, error) {
	return login(base, username, password)
}

func (scraper Scraper) Navigate(collector *colly.Collector, base, url string) (*colly.Collector, *goquery.Selection, error) {
	return navigate(collector, base, url)
}

func (scraper Scraper) Post(collector *colly.Collector, base, url string, formData map[string]string) (*colly.Collector, *goquery.Selection, error) {
	return post(collector, base, url, formData)
}

func NewScraper() *Scraper {
	return &Scraper{}
}
