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
	Navigate(collector *colly.Collector, url, endpoint string) (*colly.Collector, *goquery.Selection, error)
	Post(collector *colly.Collector, url, endpoint string, formData map[string]string) (*colly.Collector, *goquery.Selection, error)
}

type ValidationProvider interface {
	Struct(s interface{}) error
}

type QuerierProvider interface {
	GetClasswork(collector *colly.Collector, params models.ClassworkRequestBody) ([]models.Classwork, error)
	GetIPRAll(collector *colly.Collector, params models.IprAllRequestBody) ([]models.IPR, error)
	GetIPR(collector *colly.Collector, params models.IprRequestBody) ([]models.IPR, error)
	GetLogin(collector *colly.Collector, params models.LoginRequestBody) ([]models.Login, error)
	GetReportCard(collector *colly.Collector, params models.ReportCardRequestBody) ([]models.ReportCard, error)
	GetSchedule(collector *colly.Collector, params models.ScheduleRequestBody) ([]models.Schedule, error)
	GetTranscript(collector *colly.Collector, params models.TranscriptRequestBody) ([]models.Transcript, error)
}

type ParserProvider interface {
	ParseClasswork(html *goquery.Selection) models.Classwork
	ParseIPR(html *goquery.Selection) models.IPR
	ParseReportCard(html *goquery.Selection) models.ReportCard
	ParseSchedule(html *goquery.Selection) models.Schedule
	ParseTranscript(html *goquery.Selection) models.Transcript
}

type Server struct {
	App       *fiber.App
	Cache     CacheProvider
	Scraper   ScraperProvider
	Validator ValidationProvider
	Querier   QuerierProvider
	Parser    ParserProvider
}
