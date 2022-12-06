package repository

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/Threqt1/HACApi/app/models"
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

type QueryProvider interface {
	GetClasswork(collector *colly.Collector, params *models.ClassworkRequestBody) ([]models.Classwork, error)
	GetIPRAll(collector *colly.Collector, params *models.IprAllRequestBody) ([]models.IPR, error)
	GetIPR(collector *colly.Collector, params *models.IprRequestBody) ([]models.IPR, error)
	GetLogin(collector *colly.Collector, params *models.LoginRequestBody) ([]models.Login, error)
	GetReportCard(collector *colly.Collector, params *models.ReportCardRequestBody) ([]models.ReportCard, error)
	GetSchedule(collector *colly.Collector, params *models.ScheduleRequestBody) ([]models.Schedule, error)
	GetTranscript(collector *colly.Collector, params *models.TranscriptRequestBody) ([]models.Transcript, error)
}

type Server struct {
	App       *fiber.App
	Cache     CacheProvider
	Scraper   ScraperProvider
	Validator ValidationProvider
	Queries   QueryProvider
}
