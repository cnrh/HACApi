package repository

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
	"github.com/gofiber/fiber/v2"
)

type CacheProvider interface {
	GetOrLogin(key string) (*colly.Collector, error)
}

type ScraperProvider interface {
	Login(base, username, password string) (*colly.Collector, error)
	Navigate(collector *colly.Collector, base, url string) (*colly.Collector, *goquery.Selection, error)
	Post(collector *colly.Collector, base, url string, formData map[string]string) (*colly.Collector, *goquery.Selection, error)
}

type ValidationProvider interface {
	Struct(s interface{}) error
}

type Server struct {
	App       *fiber.App
	Cache     CacheProvider
	Scraper   ScraperProvider
	Validator ValidationProvider
}
